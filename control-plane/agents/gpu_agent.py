import os
import sys
import time
import glob

# Standardize path resolution
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from utils.paths import STATE_DIR
from utils.state import save_json
from utils.guards import wrap_agent
from utils.logger import get_logger

logger = get_logger("gpu_agent")

GPU_JSON = os.path.join(STATE_DIR, "gpu.json")

def get_amd_stats():
    """Probe /sys for AMD GPU statistics with multiple fallbacks."""
    stats = {
        "name": "AMD Radeon HD 5770",
        "temp": None,
        "load": 0,
        "mem_used": 0,
        "mem_total": 1024,
        "active": False
    }
    
    try:
        # Fallback 1: Direct DRM Card Probe
        cards = glob.glob("/sys/class/drm/card*")
        for card in cards:
            device_path = os.path.join(card, "device")
            if not os.path.exists(device_path): continue
            
            # Check for AMD Vendor (0x1002)
            is_amd = False
            vendor_path = os.path.join(device_path, "vendor")
            if os.path.exists(vendor_path):
                with open(vendor_path, 'r') as f:
                    if "0x1002" in f.read(): is_amd = True
            
            if not is_amd: continue
            
            stats["active"] = True
            
            # Temperature extraction (Primary)
            hwmon_paths = glob.glob(os.path.join(device_path, "hwmon/hwmon*/temp1_input"))
            if hwmon_paths:
                with open(hwmon_paths[0], 'r') as f:
                    stats["temp"] = int(f.read().strip()) / 1000.0
            
            # Load/Busy status (More paths for older cards)
            busy_paths = [
                os.path.join(device_path, "gpu_busy_percent"),
                os.path.join(device_path, "utilization"),
                os.path.join(device_path, "hwmon/hwmon*/device/gpu_busy_percent")
            ]
            for bp_pattern in busy_paths:
                for bp in glob.glob(bp_pattern):
                    if os.path.exists(bp):
                        with open(bp, 'r') as f:
                            stats["load"] = int(f.read().strip())
                            break
            
            # Memory usage (Older cards use visible_vram or different sysfs entries)
            mem_keys = {
                "mem_info_vram_used": "mem_used",
                "mem_info_vram_total": "mem_total",
                "mem_info_vis_vram_used": "mem_used",
                "vram_visible": "mem_total"
            }
            for sys_file, key in mem_keys.items():
                p = os.path.join(device_path, sys_file)
                if os.path.exists(p):
                    try:
                        with open(p, 'r') as f:
                            val = int(f.read().strip())
                            # Some files report bytes, some report MB. 5770 is ~1024MB.
                            if val > 10000: # Bytes or KB
                                stats[key] = int(val / (1024 * 1024))
                            else: # MB
                                stats[key] = val
                    except: pass
            
            # Final sanity check for VRAM total
            if stats["mem_total"] < 128: stats["mem_total"] = 1024
            
            return stats

        # Fallback 2: General HWMON Probe
        hwmon_list = glob.glob("/sys/class/hwmon/hwmon*")
        for h in hwmon_list:
            name_path = os.path.join(h, "name")
            if os.path.exists(name_path):
                with open(name_path, 'r') as f:
                    name = f.read().strip()
                    if name in ["radeon", "amdgpu", "nouveau"]:
                        stats["active"] = True
                        temp_path = os.path.join(h, "temp1_input")
                        if os.path.exists(temp_path):
                            with open(temp_path, 'r') as f:
                                stats["temp"] = int(f.read().strip()) / 1000.0
                        return stats

    except Exception as e:
        logger.error(f"GPU probe failed: {e}")
        
    return stats

def run_tick():
    """Single tick for the agent wrapper."""
    stats = get_amd_stats()
    save_json(GPU_JSON, stats)
    
    if stats["active"]:
        logger.info(f"[GPU] {stats['name']} detected: {stats['temp']}°C")

if __name__ == "__main__":
    wrap_agent("gpu_agent", run_tick, interval=10)
