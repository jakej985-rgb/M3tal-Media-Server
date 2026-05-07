import subprocess
import json
import sys
import os
import time
import collections
try:
    import psutil
except ImportError:
    psutil = None

if psutil:
    # Seed the CPU measurement (Audit fix 2.1)
    psutil.cpu_percent(interval=None)

# Add current dir to path for utils
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from utils.paths import METRICS_JSON, STATE_DIR
from utils.state import save_json, safe_replace
from utils.guards import wrap_agent
from utils.logger import get_logger

logger = get_logger("metrics")
HISTORY_CSV = os.path.join(STATE_DIR, "metrics-history.csv")
MAX_HISTORY_ENTRIES = 5000 
_last_docker_error_log = 0
_last_net_io = None
_last_net_time = None

def get_network_metrics():
    """Calculate host network throughput (MB/s) and load % (Audit fix 11.2)"""
    global _last_net_io, _last_net_time
    now = time.time()
    
    metrics = {"down": 0.0, "up": 0.0, "load": 0.0}
    
    # 🕵️ Host-Aware Network Monitoring (V6.7)
    # If we are in a container, psutil.net_io_counters() only sees container traffic.
    # We must read /host/proc/net/dev directly for host-level stats.
    host_net_dev = Path("/host/proc/net/dev")
    
    try:
        current_recv = 0
        current_sent = 0
        
        if host_net_dev.exists():
            with open(host_net_dev, "r") as f:
                lines = f.readlines()
                for line in lines[2:]: # Skip header lines
                    parts = line.split()
                    if len(parts) > 9:
                        # Index 1: Receive Bytes, Index 9: Transmit Bytes
                        current_recv += int(parts[1])
                        current_sent += int(parts[9])
        elif psutil:
            io = psutil.net_io_counters()
            current_recv = io.bytes_recv
            current_sent = io.bytes_sent
        else:
            return metrics

        if _last_net_io is not None:
            dt = now - _last_net_time
            if dt > 0.1:
                bytes_recv = current_recv - _last_net_io["recv"]
                bytes_sent = current_sent - _last_net_io["sent"]
                
                metrics["down"] = round((max(0, bytes_recv) / dt) / (1024*1024), 2)
                metrics["up"] = round((max(0, bytes_sent) / dt) / (1024*1024), 2)
                
                capacity_mbs = 125.0
                total_mbs = metrics["down"] + metrics["up"]
                metrics["load"] = round(min(100.0, (total_mbs / capacity_mbs) * 100), 1)

        _last_net_io = {"recv": current_recv, "sent": current_sent}
        _last_net_time = now
    except Exception as e:
        logger.error(f"Failed to get network metrics: {e}")
        
    return metrics

def get_system_metrics():
    metrics = {"cpu": 0.0, "mem": 0.0, "timestamp": int(time.time())}
    if psutil:
        try:
            metrics["cpu"] = psutil.cpu_percent(interval=None)
            metrics["mem"] = psutil.virtual_memory().percent
            mem_gb = psutil.virtual_memory().used / (1024**3)
            metrics["mem_gb"] = round(mem_gb, 1)
            logger.info(f"[DIAG] System metrics: CPU={metrics['cpu']}% MEM={metrics['mem']}% MEM_GB={metrics.get('mem_gb', 0)}GB")
        except Exception as e:
            logger.error(f"Failed to get psutil metrics: {e}")
    else:
        logger.warning("[DIAG] psutil not available — using fallback")
        # Fallback if psutil not available (Audit fix 2.1 - removing shell=True)
        try:
            if sys.platform != "win32":
                # Manual /proc reads are safer than shell=True + grep/awk
                if os.path.exists("/proc/stat"):
                    # Audit Fix M1: Basic /proc/stat CPU fallback
                    with open("/proc/stat", "r") as f:
                        line = f.readline()
                    parts = line.split()
                    if len(parts) >= 5:
                        idle = int(parts[4])
                        total = sum(int(p) for p in parts[1:])
                        metrics["cpu"] = round(100.0 * (1.0 - idle / total), 1)
        except Exception as e:
            logger.debug(f"Fallback metrics failed: {e}")
    return metrics

_limit_cache = {}
_last_limit_refresh = 0

def get_container_limits():
    """Fetch and cache container limits to avoid hammering the API (Audit 4.8)"""
    global _limit_cache, _last_limit_refresh
    now = time.time()
    if now - _last_limit_refresh < 60 and _limit_cache:
        return _limit_cache

    limits = {}
    try:
        # Fetch all container limits in one go
        cmd = ["docker", "inspect", "--format", "{{.Name}} {{.HostConfig.Memory}} {{.HostConfig.NanoCpus}}"]
        # Filter for running containers to keep it fast
        cmd_running = ["docker", "ps", "--format", "{{.Names}}"]
        running = subprocess.run(cmd_running, capture_output=True, text=True).stdout.strip().split('\n')
        
        proc = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
        if proc.returncode == 0:
            for line in proc.stdout.strip().split('\n'):
                parts = line.split()
                if len(parts) >= 3:
                    name = parts[0].lstrip('/')
                    if name in running or f"m3tal-{name}" in running:
                        mem = int(parts[1])
                        nano_cpus = int(parts[2])
                        limits[name] = {
                            "mem_limit": mem if mem > 0 else 0,
                            "cpu_limit": round(nano_cpus / 1e9, 2) if nano_cpus > 0 else 0
                        }
        _limit_cache = limits
        _last_limit_refresh = now
        logger.info(f"[DIAG] Discovered limits for {len(limits)} containers")
    except Exception as e:
        logger.debug(f"Limit discovery failed: {e}")
    return limits

def get_container_metrics():
    container_stats = []
    limits = get_container_limits()
    docker_host = os.environ.get("DOCKER_HOST", "(not set)")
    try:
        cmd = ["docker", "stats", "--no-stream", "--format", "{{json .}}"]
        result = subprocess.run(cmd, capture_output=True, text=True, check=True, timeout=30)
        for line in result.stdout.strip().split('\n'):
            if line:
                try:
                    raw = json.loads(line)
                    name = raw.get("Name")
                    cpu_str = str(raw.get("CPUPerc", "0")).replace("%", "")
                    mem_str = str(raw.get("MemPerc", "0")).replace("%", "")
                    
                    cpu_val = float(cpu_str) if cpu_str else 0.0
                    mem_val = float(mem_str) if mem_str else 0.0
                    
                    limit = limits.get(name) or {}
                    cpu_limit = limit.get("cpu_limit", 0)
                    mem_limit_bytes = limit.get("mem_limit", 0)
                    
                    if mem_limit_bytes > 0:
                        logger.debug(f"[LIMIT] {name}: CPU={cpu_limit}, MEM={mem_limit_bytes}")

                    # Calculate Pressure-relative CPU %
                    cpu_pressure = 0
                    if cpu_limit > 0:
                        cpu_pressure = round((cpu_val / (cpu_limit * 100)) * 100, 1)
                    else:
                        cpu_pressure = cpu_val

                    container_stats.append({
                        "name": name,
                        "cpu": cpu_val,
                        "mem": mem_val,
                        "mem_usage": raw.get("MemUsage"),
                        "cpu_limit": cpu_limit,
                        "cpu_pressure": cpu_pressure,
                        "mem_limit_bytes": mem_limit_bytes
                    })
                except (json.JSONDecodeError, ValueError) as e:
                    continue
    except Exception as e:
        logger.error(f"[DIAG] docker stats FAILED: {e}")
    return container_stats

def append_history(system, containers):
    """Batch 6 T1: Append metrics to history CSV with rotation."""
    ts = system["timestamp"]
    new_lines = []
    
    # Pre-format lines for batch write
    new_lines.append(f"{ts},host,{system['cpu']},{system['mem']}\n")
    for c in containers:
        new_lines.append(f"{ts},{c['name']},{c['cpu']},{c['mem']}\n")
        
    header = "timestamp,name,cpu,mem\n"
    
    try:
        # Step 1: Ensure header and Append new data
        mode = 'a' if os.path.isfile(HISTORY_CSV) and os.path.getsize(HISTORY_CSV) > 0 else 'w'
        try:
            with open(HISTORY_CSV, mode) as f:
                if mode == 'w':
                    f.write(header)
                f.writelines(new_lines)
        except PermissionError:
            # Audit fix: don't crash if locked, just log and skip this cycle
            logger.warning("metrics-history.csv is currently locked. Skipping this history record.")
            return

        # Step 2: Periodic Rotation check (every 10 minutes)
        last_prune_file = os.path.join(STATE_DIR, "last_prune.json")
        last_prune_ts = 0
        try:
            if os.path.exists(last_prune_file):
                with open(last_prune_file, 'r') as pf:
                    last_prune_ts = json.loads(pf.read().strip() or '{}').get("ts", 0)
        except Exception as pe:
            logger.debug(f"Non-critical last_prune read failure: {pe}")

        if ts - last_prune_ts > 600:
            try:
                if os.path.exists(HISTORY_CSV):
                    with open(HISTORY_CSV, "r") as f:
                        # Use maxlen to automatically prune oldest line on read
                        all_lines = collections.deque(f, maxlen=MAX_HISTORY_ENTRIES)
                    
                    # Write to temp file then replace (Audit Fix 8)
                    tmp_csv = f"{HISTORY_CSV}.tmp"
                    with open(tmp_csv, "w") as f:
                        if all_lines and not all_lines[0].startswith("timestamp"):
                            f.write(header)
                        f.writelines(all_lines)
                    
                    safe_replace(tmp_csv, HISTORY_CSV)
                
                # Update last prune timestamp
                with open(last_prune_file, 'w') as pf:
                    json.dump({"ts": ts}, pf)
                logger.info("Rotated metrics-history.csv")
            except PermissionError:
                logger.warning("Rotation failed: metrics-history.csv is locked.")
            except Exception as e:
                logger.error(f"Rotation error: {e}")
    except Exception as e:
        logger.error(f"Failed to process history: {e}")

def collect_all_metrics():
    system = get_system_metrics()
    containers = get_container_metrics()
    network = get_network_metrics()
    
    data = {
        "system": system,
        "containers": containers,
        "network": network,
        "timestamp": system["timestamp"],
        "cpu": system["cpu"]
    }
    
    if save_json(METRICS_JSON, data, caller="metrics"):
        append_history(system, containers)
        logger.info(f"[DIAG] Saved metrics: CPU={system['cpu']}% MEM={system['mem']}% Containers={len(containers)}")
    else:
        logger.error(f"[DIAG] FAILED to save metrics to {METRICS_JSON}")

if __name__ == "__main__":
    wrap_agent("metrics", collect_all_metrics)
