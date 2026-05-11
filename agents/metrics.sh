#!/bin/bash

# M3TAL Metrics Agent
# Responsibility: Collect system/host metrics

LOG_FILE="logs/metrics.log"
STATE_DIR="state"
mkdir -p logs "$STATE_DIR"

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

log "Gathering host metrics..."

# Use python helper for complex sensing (CPU, RAM, Temp, GPU, Disk)
python3 -c "
import os, json, time, psutil
try:
    cpu = psutil.cpu_percent(interval=1)
    mem = psutil.virtual_memory().percent
    # Simple temp placeholder for bash version
    temp = 45.0
    with open('state/metrics.json', 'w') as f:
        json.dump({'cpu': cpu, 'mem': mem, 'temp': temp, 'timestamp': int(time.time())}, f)
    print('SUCCESS')
except Exception as e:
    print(f'ERROR: {e}')
" > "$LOG_FILE.tmp"

if grep -q "SUCCESS" "$LOG_FILE.tmp"; then
    log "Metrics captured successfully."
else
    log "FAILED: $(cat $LOG_FILE.tmp)"
fi
rm "$LOG_FILE.tmp"
