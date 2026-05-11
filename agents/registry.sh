#!/bin/bash

# M3TAL Registry Agent
# Responsibility: Track system state (Source of Truth)

LOG_FILE="logs/registry.log"
STATE_DIR="state"
mkdir -p logs "$STATE_DIR"

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

log "Updating global system registry..."

# Build the unified system.json
# We'll use python to safely merge the files
python3 -c "
import json, os, time
state = {
    'status': 'healthy',
    'last_check': time.strftime('%Y-%m-%d %H:%M:%S'),
    'timestamp': int(time.time()),
    'metrics': {},
    'containers': [],
    'anomalies': []
}

try:
    if os.path.exists('state/metrics.json'):
        with open('state/metrics.json') as f: state['metrics'] = json.load(f)
    if os.path.exists('state/containers.json'):
        with open('state/containers.json') as f: state['containers'] = json.load(f)
    if os.path.exists('state/anomalies.json'):
        with open('state/anomalies.json') as f: state['anomalies'] = json.load(f)
    
    # Calculate global status
    if len(state['anomalies']) > 0 and 'issue' in state['anomalies']:
        state['status'] = 'degraded'
        
    with open('state/system.json', 'w') as f:
        json.dump(state, f, indent=2)
    print('SUCCESS')
except Exception as e:
    print(f'ERROR: {e}')
" > "$LOG_FILE.tmp"

if grep -q "SUCCESS" "$LOG_FILE.tmp"; then
    log "Registry synchronized: state/system.json updated."
else
    log "FAILED to sync registry: $(cat $LOG_FILE.tmp)"
fi
rm "$LOG_FILE.tmp"
