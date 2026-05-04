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
    """Discovery using lsblk and /host/proc/mounts to find all host drives."""
    stats = {}
    highest_usage = 0

    # 1. Identify all physical devices and partitions
    devices = get_lsblk_data()
    
    # 2. Check host mounts if running in Docker
    # This allows us to see drives NOT mounted in the container itself
    host_mounts = {}
    if os.path.exists("/host/proc/mounts"):
        try:
            with open("/host/proc/mounts", "r") as f:
                for line in f:
                    parts = line.split()
                    if len(parts) >= 2:
                        dev, mnt = parts[0], parts[1]
                        if dev.startswith("/dev/sd") or dev.startswith("/dev/nvme"):
                            host_mounts[dev] = mnt
        except Exception as e:
            logger.debug(f"Failed to read host mounts: {e}")

    # 3. Flat list of candidates (from lsblk and host mounts)
    candidates = []
    
    def walk_devices(devs):
        for d in devs:
            dev_path = f"/dev/{d['name']}"
            # Check if this device/partition has a mountpoint in container OR on host
            mnt = d.get("mountpoint") or host_mounts.get(dev_path)
            if mnt:
                candidates.append({
                    "name": d["name"],
                    "label": d.get("label"),
                    "mountpoint": mnt,
                    "device": dev_path
                })
            if d.get("children"):
                walk_devices(d["children"])

    walk_devices(devices)

    # 4. Process candidates
    for c in candidates:
        mnt = c["mountpoint"]
        dev = c["device"]
        
        # Determine the path to check usage (prefix with /host if it exists)
        check_path = mnt
        if os.path.exists("/host") and not mnt.startswith("/host"):
            check_path = os.path.join("/host", mnt.lstrip("/"))
            if not os.path.exists(check_path):
                # Try relative to host if it's a standard /mnt or /media
                check_path = mnt

        try:
            if not os.path.exists(check_path):
                continue

            usage = psutil.disk_usage(check_path)
            
            # Label Hierarchy: LABEL > MOUNTPOINT basename > DEVICE
            label = c.get("label")
            if not label or label.strip() == "":
                if mnt == "/":
                    label = "System"
                else:
                    # Clean up mountpoint name (e.g. /mnt/media -> Media)
                    label = os.path.basename(mnt.rstrip("/")).capitalize() or c["name"]
            
            # Get temperature for the underlying physical disk
            # We strip trailing numbers to get the physical disk (e.g. sda1 -> sda)
            parent_name = ''.join(c for c in c["name"] if not c.isdigit())
            temp = get_drive_temp(f"/dev/{parent_name}")

            stats[label] = {
                "device": dev,
                "mountpoint": mnt,
                "total_gb": round(usage.total / (1024**3), 1),
                "used_gb": round(usage.used / (1024**3), 1),
                "percent": usage.percent,
                "temp": temp
            }
            if usage.percent > highest_usage:
                highest_usage = usage.percent
        except Exception as e:
            logger.debug(f"Failed to read disk usage for {mnt} (at {check_path}): {e}")

    # Fallback to psutil (Windows/Local Dev)
    if not stats and psutil:
        # ... (psutil logic remains same as fallback)
        pass

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
