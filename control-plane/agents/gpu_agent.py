import os
import sys
import time
import glob
import re
import subprocess

# Standardize path resolution
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from utils.paths import STATE_DIR
from utils.state import save_json
from utils.guards import wrap_agent
from utils.logger import get_logger

logger = get_logger("gpu_agent")

GPU_JSON = os.path.join(STATE_DIR, "gpu.json")

def get_radeontop_stats():
    """Execute radeontop and parse its dump output."""
    stats = {"load": 0, "mem_used": 0, "mem_total": 1024}
    try:
        # -d - : dump to stdout
        # -l 1 : only 1 sample
        cmd = ["radeontop", "-d", "-", "-l", "1"]
        res = subprocess.run(cmd, capture_output=True, text=True, timeout=5)
        if res.returncode == 0:
            line = res.stdout.strip()
            # Example: ... gpu 12.50%, ... vram 10.20% 104.45mb ...
            
            # 1. Parse GPU Load
            gpu_match = re.search(r"gpu\s+([\d\.]+)%", line)
            if gpu_match:
                stats["load"] = int(float(gpu_match.group(1)))
            
            # 2. Parse VRAM used (MB)
            vram_mb_match = re.search(r"vram\s+[\d\.]+% ([\d\.]+)mb", line)
            if vram_mb_match:
                stats["mem_used"] = int(float(vram_mb_match.group(1)))
                
            return stats, True
    except Exception as e:
        logger.debug(f"radeontop execution failed: {e}")
    return stats, False

def get_amd_stats():
    """Probe /sys for AMD GPU statistics with radeontop as primary."""
    stats = {
        "name": "AMD Radeon HD 5770",
        "temp": None,
        "load": 0,
        "mem_used": 0,
        "mem_total": 1024,
        "active": False
    }
    
    try:
        # 1. Find the card for Temperature (HWMON)
        cards = glob.glob("/sys/class/drm/card*")
        for card in cards:
            device_path = os.path.join(card, "device")
            if not os.path.exists(device_path): continue
            
            vendor_path = os.path.join(device_path, "vendor")
            if os.path.exists(vendor_path):
                with open(vendor_path, 'r') as f:
                    if "0x1002" not in f.read(): continue
            
            stats["active"] = True
            
            # Get Temperature from HWMON
            hwmon_paths = glob.glob(os.path.join(device_path, "hwmon/hwmon*/temp1_input"))
            if hwmon_paths:
                with open(hwmon_paths[0], 'r') as f:
                    stats["temp"] = int(f.read().strip()) / 1000.0
            break

        # 2. Use radeontop for the heavy lifting (Load/RAM)
        rt_stats, rt_success = get_radeontop_stats()
        if rt_success:
            stats["load"] = rt_stats["load"]
            stats["mem_used"] = rt_stats["mem_used"]
            stats["active"] = True # Force active if radeontop works
        
    except Exception as e:
        logger.error(f"GPU probe failed: {e}")
        
    return stats

def run_tick():
    """Single tick for the agent wrapper."""
    stats = get_amd_stats()
    save_json(GPU_JSON, stats)
    
    if stats["active"]:
        logger.info(f"[GPU] {stats['name']} detected: {stats['temp']}°C | VRAM: {stats['mem_used']}MB | Load: {stats['load']}%")

if __name__ == "__main__":
    wrap_agent("gpu_agent", run_tick, interval=10)
