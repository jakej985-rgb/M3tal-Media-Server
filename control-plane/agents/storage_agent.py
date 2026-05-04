import os
import sys
import time
import subprocess
import json

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

_smartctl_path = None
_device_types = {} # Cache for scan results

def get_drive_temp(dev_path):
    global _smartctl_path, _device_types
    if not dev_path or not dev_path.startswith('/dev/'):
        return None
    
    # Discovery logic
    if not _smartctl_path:
        import shutil
        _smartctl_path = shutil.which("smartctl") or "/usr/sbin/smartctl"
        if not os.path.exists(_smartctl_path):
            return None

    # Perform a scan once to find device types (handles USB/RAID bridges)
    if not _device_types:
        try:
            scan_res = subprocess.run([_smartctl_path, "--scan"], capture_output=True, text=True, timeout=5)
            for line in scan_res.stdout.splitlines():
                # Format: /dev/sda -d sat # ...
                if " -d " in line:
                    parts = line.split(" -d ")
                    dev = parts[0].strip()
                    d_type = parts[1].split()[0]
                    _device_types[dev] = d_type
        except: pass

    # Try cached type first, then fallback to common types
    dev_type = _device_types.get(dev_path)
    types_to_try = [dev_type] if dev_type else [None, "sat", "nvme", "scsi"]
    
    for d_type in types_to_try:
        try:
            cmd = [_smartctl_path, "-a"]
            if d_type: cmd += ["-d", d_type]
            cmd.append(dev_path)
                
            res = subprocess.run(cmd, capture_output=True, text=True, timeout=4)
            
            import re
            patterns = [
                # Catch the raw value at the end of the line for standard attributes (190, 194, etc)
                r"(?:Temperature_Celsius|Airflow_Temperature_Cel|Temperature_Internal).*?\s+(\d+)(?:\s+\(.*\))?$",
                # Catch NVMe or simplified outputs
                r"(?:Composite\s+Temperature|Temperature:).*?(\d+)",
                r"Current\s+Drive\s+Temperature:\s+(\d+)",
            ]
            
            for pattern in patterns:
                # Use MULTILINE to ensure $ works per line
                match = re.search(pattern, res.stdout, re.IGNORECASE | re.MULTILINE)
                if match:
                    val = float(match.group(1))
                    if 5 < val < 100: return val # Sanity check for real temps
            
            # If SMART is disabled, try to enable it once
            if "SMART support is: Disabled" in res.stdout:
                subprocess.run([_smartctl_path, "-s", "on", dev_path], capture_output=True, timeout=2)
                
        except: pass
            
    return None
            
    return None

def get_lsblk_data():
    """Execute lsblk and return JSON data."""
    try:
        output = subprocess.check_output(
            "lsblk -J -o NAME,LABEL,MOUNTPOINT,TYPE,SIZE",
            shell=True
        ).decode()
        return json.loads(output).get("blockdevices", [])
    except Exception as e:
        logger.debug(f"lsblk failed: {e}")
        return []

def get_free_space(mount):
    """Get free space using df -BG."""
    try:
        # If in Docker, we must check /host mount
        check_path = mount
        if os.path.exists("/host") and not mount.startswith("/host"):
            check_path = os.path.join("/host", mount.lstrip("/"))
        
        output = subprocess.check_output(
            f"df -BG {check_path} | tail -1",
            shell=True
        ).decode()
        # Output format: Filesystem 1G-blocks Used Available Use% MountedOn
        parts = output.split()
        if len(parts) >= 4:
            return parts[3].replace("G", "") # Return just the number
    except:
        pass
    return "N/A"

def get_disk_stats():
    """Enumerate physical disks and map to partitions."""
    stats = {}
    highest_usage = 0
    
    blocks = get_lsblk_data()
    
    for block in blocks:
        if block.get("type") != "disk":
            continue
            
        disk_name = block["name"]
        
        # Look for children (partitions)
        if "children" in block:
            for part in block["children"]:
                mount = part.get("mountpoint")
                
                # If not mounted in container, check if it's mounted on host
                # We can check /host/proc/mounts to see what the host sees
                if not mount and os.path.exists("/host/proc/mounts"):
                    try:
                        with open("/host/proc/mounts", "r") as f:
                            for line in f:
                                if f"/dev/{part['name']}" in line:
                                    mount = line.split()[1]
                                    break
                    except: pass
                
                if mount:
                    label = part.get("label")
                    if not label or label.strip() == "":
                        label = mount.split("/")[-1].capitalize() or part["name"]
                    
                    if label == "" or label == "/": label = "System"

                    free = get_free_space(mount)
                    # Use the disk_name (e.g. sda) for the temperature sensor
                    temp = get_drive_temp(f"/dev/{disk_name}")
                    
                    logger.info(f"[STORAGE] Found {label} on {disk_name} ({mount}): {free}G, {temp}°C")
                    
                    # Get usage % for the icon color
                    try:
                        check_path = mount
                        if os.path.exists("/host") and not mount.startswith("/host"):
                            check_path = os.path.join("/host", mount.lstrip("/"))
                        usage = psutil.disk_usage(check_path)
                        percent = usage.percent
                    except: percent = 0

                    stats[label] = {
                        "free": free,
                        "temp": temp,
                        "percent": percent
                    }
                    if percent > highest_usage: highest_usage = percent

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
