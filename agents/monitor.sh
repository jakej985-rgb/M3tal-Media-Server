#!/bin/bash

# M3TAL Monitor Agent
# Responsibility: Detect container/service state

LOG_FILE="logs/monitor.log"
STATE_FILE="state/containers.json"
mkdir -p logs state

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

log "Sensing container fleet..."

# Use Docker CLI to get status of all containers
# Format: Name,Status,State,Labels
docker ps -a --format '{"name":"{{.Names}}","status":"{{.Status}}","state":"{{.State}}","labels":"{{.Labels}}"}' > "$STATE_FILE.tmp"

# Wrap in JSON array
(echo "["; sed 's/$/,/' "$STATE_FILE.tmp" | sed '$ s/,$//'; echo "]") > "$STATE_FILE"
rm "$STATE_FILE.tmp"

log "Detected $(grep -c "name" "$STATE_FILE") containers."
log "Monitor completed successfully."
