import os
import sys
import time

try:
    import psutil
except ImportError:
    psutil = None

# Standardize path resolution
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from utils.paths import STORAGE_JSON, DATA_DIR
from utils.state import save_json
from utils.guards import wrap_agent
from utils.logger import get_logger

logger = get_logger("storage_agent")

import subprocess
import json

def get_lsblk_data():
    try:
        # We need -b for bytes if we want exact numbers, or just rely on psutil for usage
        # TRAN (transport) helps identify USB/SATA/NVME
        res = subprocess.run(
            ["lsblk", "-J", "-o", "NAME,LABEL,MOUNTPOINT,SIZE,FSTYPE,TYPE,UUID,TRAN"],
            capture_output=True, text=True, timeout=5
        )
        if res.returncode == 0:
            return json.loads(res.stdout).get("blockdevices", [])
    except Exception as e:
        logger.debug(f"lsblk failed: {e}")
    return []

def get_drive_temp(dev_path):
    if not dev_path or not dev_path.startswith('/dev/'):
        return None
    try:
        # 🛡️ SECURITY: smartctl requires root but we are running as root in runtime
        res = subprocess.run(
            ["smartctl", "-A", dev_path],
            capture_output=True, text=True, timeout=2
        )
        # Standard SMART output parsing for SATA/SAS
        for line in res.stdout.splitlines():
            if "Temperature_Celsius" in line:
                parts = line.split()
                if len(parts) >= 10:
                    return float(parts[9])
            # NVMe uses a different format
            if "Temperature:" in line:
                parts = line.split()
                if len(parts) >= 2:
                    # Remove 'Celsius' or non-numeric chars
                    temp_str = ''.join(c for c in parts[1] if c.isdigit() or c == '.')
                    return float(temp_str)
    except Exception:
        pass
    return None

def get_disk_stats():
    """Discovery using lsblk (Linux) with fallback to psutil (Windows/Basic)."""
    stats = {}
    highest_usage = 0

    # 1. Try lsblk discovery (Linux specific, works better in Docker if /dev is exposed)
    devices = get_lsblk_data()
    
    # Flat list of partitions that have mountpoints
    found_partitions = []
    
    def walk_devices(devs):
        for d in devs:
            if d.get("mountpoint"):
                found_partitions.append(d)
            if d.get("children"):
                walk_devices(d["children"])

    walk_devices(devices)

    # 2. Process found partitions
    for p in found_partitions:
        mnt = p["mountpoint"]
        dev = f"/dev/{p['name']}"
        
        try:
            usage = psutil.disk_usage(mnt)
            
            # Label Hierarchy: LABEL > MOUNTPOINT basename > DEVICE
            label = p.get("label")
            if not label or label.strip() == "":
                if mnt == "/":
                    label = "System"
                else:
                    label = os.path.basename(mnt.rstrip("/")) or p["name"]
            
            # Get temperature for the underlying physical disk
            # We strip trailing numbers to get the physical disk (e.g. sda1 -> sda)
            parent_dev = "/dev/" + ''.join(c for c in p["name"] if not c.isdigit())
            temp = get_drive_temp(parent_dev)

            stats[label] = {
                "device": dev,
                "mountpoint": mnt,
                "fstype": p.get("fstype", "unknown"),
                "total_gb": round(usage.total / (1024**3), 1),
                "used_gb": round(usage.used / (1024**3), 1),
                "free_gb": round(usage.free / (1024**3), 1),
                "percent": usage.percent,
                "temp": temp
            }
            if usage.percent > highest_usage:
                highest_usage = usage.percent
        except Exception as e:
            logger.debug(f"Failed to read disk usage for {mnt}: {e}")

    # 3. Fallback to psutil if no drives found (e.g. Windows dev)
    if not stats and psutil:
        try:
            partitions = psutil.disk_partitions(all=False)
            for p in partitions:
                if 'cdrom' in p.opts or p.fstype == '': continue
                try:
                    usage = psutil.disk_usage(p.mountpoint)
                    label = p.mountpoint if p.mountpoint != "/" else "System"
                    stats[label] = {
                        "device": p.device,
                        "mountpoint": p.mountpoint,
                        "fstype": p.fstype,
                        "total_gb": round(usage.total / (1024**3), 1),
                        "used_gb": round(usage.used / (1024**3), 1),
                        "free_gb": round(usage.free / (1024**3), 1),
                        "percent": usage.percent,
                        "temp": None
                    }
                    if usage.percent > highest_usage:
                        highest_usage = usage.percent
                except Exception: continue
        except Exception as e:
            logger.error(f"Fallback partition listing failed: {e}")

    return stats, highest_usage

def get_io_stats():
    if not psutil:
        return None
    try:
        io = psutil.disk_io_counters()
        if io:
            return {
                "read_count": io.read_count,
                "write_count": io.write_count,
                "read_bytes": io.read_bytes,
                "write_bytes": io.write_bytes
            }
    except Exception as e:
        logger.debug(f"Failed to read disk IO: {e}")
    return None

def collect_storage():
    logger.info("[STORAGE] Collecting storage metrics...")
    
    disk_stats, highest_usage = get_disk_stats() or ({}, 0)
    io_stats = get_io_stats()

    data = {
        "disks": disk_stats,
        "io": io_stats,
        "timestamp": int(time.time()),
        "status": "healthy"
    }

    # Evaluate Thresholds
    if highest_usage >= 95:
        data["status"] = "critical"
        logger.warning(f"[STORAGE] CRITICAL Disk Usage detected: {highest_usage}%")
    elif highest_usage >= 85:
        data["status"] = "warning"
        logger.info(f"[STORAGE] Warning Disk Usage detected: {highest_usage}%")

    save_json(STORAGE_JSON, data, caller="storage_agent")
    logger.info(f"[STORAGE] Saved metrics. Highest usage: {highest_usage}%")

if __name__ == "__main__":
    wrap_agent("storage_agent", collect_storage)
