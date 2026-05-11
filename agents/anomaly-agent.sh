#!/bin/bash

# M3TAL Anomaly Agent
# Responsibility: Detect abnormal patterns

LOG_FILE="logs/anomaly.log"
STATE_DIR="state"
mkdir -p logs "$STATE_DIR"

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

log "Analyzing container fleet for anomalies..."

# Check containers.json for crashed services
if [ -f "state/containers.json" ]; then
    CRASHED=$(grep -E "exited|dead" state/containers.json | wc -l)
    if [ "$CRASHED" -gt 0 ]; then
        log "CRITICAL: Detected $CRASHED crashed containers!"
        echo "{\"issue\": \"container_crash\", \"count\": $CRASHED, \"timestamp\": $(date +%s)}" > state/anomalies.json
    else
        log "System health nominal. No container anomalies."
        echo "[]" > state/anomalies.json
    fi
else
    log "WARN: No container state found. Skipping check."
fi

log "Anomaly analysis completed."
