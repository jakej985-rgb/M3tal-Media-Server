import os
import sys
import time
import glob
import re

# Standardize path resolution
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from utils.paths import STATE_DIR
from utils.state import save_json
from utils.guards import wrap_agent
from utils.logger import get_logger

logger = get_logger("gpu_agent")

GPU_JSON = os.path.join(STATE_DIR, "gpu.json")

def parse_pm_info(path):
    """Parse /sys/kernel/debug/dri/*/radeon_pm_info for load/clocks."""
    data = {"load": 0, "temp": None}
    try:
        if os.path.exists(path):
            with open(path, 'r') as f:
                content = f.read()
                # Look for "GPU load: X%" or similar patterns
                load_match = re.search(r"GPU load:\s*(\d+)%", content)
                if load_match:
                    data["load"] = int(load_match.group(1))
                
                # Temperature fallback if not found in hwmon
                temp_match = re.search(r"temperature:\s*(\d+)", content)
                if temp_match:
                    data["temp"] = float(temp_match.group(1))
    except: pass
    return data

def get_amd_stats():
    """Probe /sys for AMD GPU statistics with Deep Probe fallbacks."""
    stats = {
        "name": "AMD Radeon HD 5770",
        "temp": None,
        "load": 0,
        "mem_used": 0,
        "mem_total": 1024,
        "active": False
    }
    
    try:
        # 1. Find the card and device path
        cards = glob.glob("/sys/class/drm/card*")
        for card in cards:
            device_path = os.path.join(card, "device")
            if not os.path.exists(device_path): continue
            
            # Check for AMD Vendor (0x1002)
            vendor_path = os.path.join(device_path, "vendor")
            if os.path.exists(vendor_path):
                with open(vendor_path, 'r') as f:
                    if "0x1002" not in f.read(): continue
            
            stats["active"] = True
            
            # 2. Temperature (HWMON)
            hwmon_paths = glob.glob(os.path.join(device_path, "hwmon/hwmon*/temp1_input"))
            if hwmon_paths:
                with open(hwmon_paths[0], 'r') as f:
                    stats["temp"] = int(f.read().strip()) / 1000.0
            
            # 3. Memory (Standard & Visible paths)
            vram_paths = [
                ("mem_info_vram_used", "mem_used"),
                ("mem_info_vram_total", "mem_total"),
                ("mem_info_vis_vram_used", "mem_used"),
                ("vram_visible", "mem_total")
            ]
            for sys_file, key in vram_paths:
                p = os.path.join(device_path, sys_file)
                if os.path.exists(p):
                    try:
                        with open(p, 'r') as f:
                            val = int(f.read().strip())
                            if val > 10000: stats[key] = int(val / (1024 * 1024))
                            else: stats[key] = val
                    except: pass

            # 4. Deep Probe: debugfs (radeon_pm_info)
            # This is the gold mine for older Juniper XT cards
            card_idx = card.replace("/sys/class/drm/card", "")
            debug_path = f"/sys/kernel/debug/dri/{card_idx}/radeon_pm_info"
            pm_data = parse_pm_info(debug_path)
            
            if pm_data["load"] > 0:
                stats["load"] = pm_data["load"]
            if stats["temp"] is None:
                stats["temp"] = pm_data["temp"]

            # Fallback for Load: Check gpu_busy_percent
            busy_path = os.path.join(device_path, "gpu_busy_percent")
            if os.path.exists(busy_path) and stats["load"] == 0:
                with open(busy_path, 'r') as f:
                    stats["load"] = int(f.read().strip())

            break 

        # Final sanity check for VRAM
        if stats["mem_total"] < 128: stats["mem_total"] = 1024
        # If VRAM used is still 0, check the 'vram_usage' file if it exists
        if stats["active"] and stats["mem_used"] == 0:
            usage_path = os.path.join(device_path, "vram_usage")
            if os.path.exists(usage_path):
                with open(usage_path, 'r') as f:
                    stats["mem_used"] = int(int(f.read().strip()) / (1024 * 1024))

    except Exception as e:
        logger.error(f"GPU probe failed: {e}")
        
    return stats

def run_tick():
    """Single tick for the agent wrapper."""
    stats = get_amd_stats()
    save_json(GPU_JSON, stats)
    
    if stats["active"]:
        logger.info(f"[GPU] {stats['name']} detected: {stats['temp']}°C | Load: {stats['load']}%")

if __name__ == "__main__":
    wrap_agent("gpu_agent", run_tick, interval=10)
