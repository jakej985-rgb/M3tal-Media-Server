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

def get_amd_stats():
    """Probe /sys for AMD GPU statistics with exhaustive VRAM search."""
    stats = {
        "name": "AMD Radeon HD 5770",
        "temp": None,
        "load": 0,
        "mem_used": 0,
        "mem_total": 1024,
        "active": False
    }
    
    try:
        # 1. Find the card
        cards = glob.glob("/sys/class/drm/card*")
        for card in cards:
            device_path = os.path.join(card, "device")
            if not os.path.exists(device_path): continue
            
            # Vendor Check
            vendor_path = os.path.join(device_path, "vendor")
            if os.path.exists(vendor_path):
                with open(vendor_path, 'r') as f:
                    if "0x1002" not in f.read(): continue
            
            stats["active"] = True
            card_idx = card.replace("/sys/class/drm/card", "")
            
            # 2. Temperature (Priority: Chip Sensor)
            hwmon_paths = glob.glob(os.path.join(device_path, "hwmon/hwmon*/temp1_input"))
            if hwmon_paths:
                with open(hwmon_paths[0], 'r') as f:
                    stats["temp"] = int(f.read().strip()) / 1000.0
            
            # 3. Exhaustive Memory Search
            # Try debugfs first (most reliable for old radeon cards)
            vram_usage_debug = f"/sys/kernel/debug/dri/{card_idx}/radeon_vram_usage"
            if os.path.exists(vram_usage_debug):
                try:
                    with open(vram_usage_debug, 'r') as f:
                        # Format: "vram_usage: 12345678" (bytes)
                        content = f.read().strip()
                        val = re.search(r"(\d+)", content)
                        if val:
                            stats["mem_used"] = int(int(val.group(1)) / (1024 * 1024))
                except: pass

            # Fallback to standard sysfs paths
            mem_paths = [
                "mem_info_vram_used", "mem_info_vis_vram_used", 
                "vram_usage", "vram_visible"
            ]
            for mp in mem_paths:
                if stats["mem_used"] > 0: break
                p = os.path.join(device_path, mp)
                if os.path.exists(p):
                    try:
                        with open(p, 'r') as f:
                            val = int(f.read().strip())
                            if val > 10000: stats["mem_used"] = int(val / (1024 * 1024))
                            else: stats["mem_used"] = val
                    except: pass

            # 4. Load Detection (Power Level Fallback)
            # Old cards don't have "load %", so we check if the power state is high
            pm_info = f"/sys/kernel/debug/dri/{card_idx}/radeon_pm_info"
            if os.path.exists(pm_info):
                try:
                    with open(pm_info, 'r') as f:
                        content = f.read()
                        # If we see high clocks, assume some load
                        if "power level 2" in content or "power level 1" in content:
                            stats["load"] = 5 if "level 1" in content else 15
                except: pass
            
            # Also check gpu_busy_percent if it exists
            busy_p = os.path.join(device_path, "gpu_busy_percent")
            if os.path.exists(busy_p):
                with open(busy_p, 'r') as f:
                    stats["load"] = int(f.read().strip())

            break 

        # Sanity check
        if stats["mem_total"] < 128: stats["mem_total"] = 1024
        
    except Exception as e:
        logger.error(f"GPU probe failed: {e}")
        
    return stats

def run_tick():
    """Single tick for the agent wrapper."""
    stats = get_amd_stats()
    save_json(GPU_JSON, stats)
    
    if stats["active"]:
        logger.info(f"[GPU] {stats['name']} detected: {stats['temp']}°C | VRAM: {stats['mem_used']}MB")

if __name__ == "__main__":
    wrap_agent("gpu_agent", run_tick, interval=10)
