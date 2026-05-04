import os
import sys
import time
import glob

# Standardize path resolution
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from utils.paths import DATA_DIR
from utils.state import save_json
from utils.guards import wrap_agent
from utils.logger import get_logger

logger = get_logger("gpu_agent")

GPU_JSON = os.path.join(DATA_DIR, "gpu.json")

def get_amd_stats():
    """Probe /sys for AMD GPU statistics."""
    stats = {
        "name": "AMD Radeon HD 5770",
        "temp": None,
        "load": 0,
        "mem_used": 0,
        "mem_total": 1024, # Default for 5770
        "active": False
    }
    
    try:
        # 1. Find the card (usually card0 or card1)
        cards = glob.glob("/sys/class/drm/card*")
        for card in cards:
            device_path = os.path.join(card, "device")
            
            # Check if it's actually an AMD card (Vendor 0x1002)
            vendor_path = os.path.join(device_path, "vendor")
            if os.path.exists(vendor_path):
                with open(vendor_path, 'r') as f:
                    if "0x1002" not in f.read():
                        continue
            
            stats["active"] = True
            
            # 2. Get Temperature
            hwmon_paths = glob.glob(os.path.join(device_path, "hwmon/hwmon*/temp1_input"))
            if hwmon_paths:
                with open(hwmon_paths[0], 'r') as f:
                    stats["temp"] = int(f.read().strip()) / 1000.0
            
            # 3. Get Load (if available on this driver)
            busy_path = os.path.join(device_path, "gpu_busy_percent")
            if os.path.exists(busy_path):
                with open(busy_path, 'r') as f:
                    stats["load"] = int(f.read().strip())
            
            # 4. Get Memory
            vram_used_path = os.path.join(device_path, "mem_info_vram_used")
            vram_total_path = os.path.join(device_path, "mem_info_vram_total")
            
            if os.path.exists(vram_used_path):
                with open(vram_used_path, 'r') as f:
                    stats["mem_used"] = int(int(f.read().strip()) / (1024 * 1024))
            if os.path.exists(vram_total_path):
                with open(vram_total_path, 'r') as f:
                    stats["mem_total"] = int(int(f.read().strip()) / (1024 * 1024))
            
            break # Found our primary AMD card
            
    except Exception as e:
        logger.debug(f"GPU probe failed: {e}")
        
    return stats

@wrap_agent("gpu_agent")
def main():
    while True:
        stats = get_amd_stats()
        save_json(GPU_JSON, stats)
        
        if stats["active"]:
            logger.debug(f"[GPU] {stats['name']}: {stats['temp']}°C, {stats['load']}% Load")
        
        time.sleep(10)

if __name__ == "__main__":
    main()
