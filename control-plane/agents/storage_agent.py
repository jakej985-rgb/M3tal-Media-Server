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

def get_drive_temp(dev_path):
    global _smartctl_path
    if not dev_path or not dev_path.startswith('/dev/'):
        return None
    
    # Discovery logic (runs once or until found)
    if not _smartctl_path:
        import shutil
        _smartctl_path = shutil.which("smartctl")
        
        if not _smartctl_path:
            # Deep search in common binary locations
            search_paths = [
                "/usr/sbin/smartctl", "/usr/bin/smartctl", "/sbin/smartctl", 
                "/bin/smartctl", "/usr/local/sbin/smartctl", "/usr/local/bin/smartctl"
            ]
            for p in search_paths:
                if os.path.exists(p):
                    _smartctl_path = p
                    break
        
        if not _smartctl_path:
            logger.info(f"[STORAGE] CRITICAL: smartctl not found in PATH or common dirs.")
            logger.info(f"[STORAGE] System PATH: {os.environ.get('PATH')}")
            return None
        else:
            logger.info(f"[STORAGE] Using smartctl at: {_smartctl_path}")

    # Try different device types if default fails (common for USB/NVMe bridges)
    types_to_try = [None, "sat", "nvme"]
    
    for d_type in types_to_try:
        try:
            cmd = [_smartctl_path, "-a", dev_path]
            if d_type:
                cmd = [_smartctl_path, "-d", d_type, "-a", dev_path]
                
            res = subprocess.run(cmd, capture_output=True, text=True, timeout=4)
            
            import re
            patterns = [
                r"(?:Temperature_Celsius|Temperature_Internal|Airflow_Temperature_Celsius|Composite\s+Temperature|Temperature:).*?(\d+)",
                r"Current\s+Drive\s+Temperature:\s+(\d+)",
            ]
            
            for pattern in patterns:
                match = re.search(pattern, res.stdout, re.IGNORECASE)
                if match:
                    return float(match.groups()[-1])
            
            # If we reached here with no match, and it's our last attempt, log why
            if d_type == types_to_try[-1] and res.stdout:
                snippet = res.stdout[:150].replace('\n', ' ')
                logger.info(f"[STORAGE] {dev_path} scan failed. Output snippet: {snippet}")
        except Exception as e:
            if d_type == types_to_try[-1]:
                logger.info(f"[STORAGE] smartctl {dev_path} error: {e}")
            pass
            
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
                    temp = get_drive_temp(f"/dev/{disk_name}")
                    
                    logger.info(f"[STORAGE] Found {label} ({mount}): {free}G, {temp}°C")
                    
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
