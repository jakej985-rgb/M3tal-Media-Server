#!/bin/bash

# M3TAL Reconcile Agent
# Responsibility: Enforce desired state (Self-healing)

LOG_FILE="logs/reconcile.log"
STATE_DIR="state"
mkdir -p logs "$STATE_DIR"

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

log "Reconciling system state..."

if [ -f "state/decisions.json" ] && grep -q "reconcile_containers" state/decisions.json; then
    log "Executing container recovery..."
    
    # Identify managed containers that are NOT running
    # We look for m3tal.managed=true and state != running
    for cid in $(docker ps -a --filter "label=m3tal.managed=true" --filter "status=exited" --filter "status=dead" -q); do
        name=$(docker inspect --format '{{.Name}}' "$cid" | sed 's/\///')
        log "ACTION: Restarting failed container: $name"
        docker restart "$cid" > /dev/null
        log "SUCCESS: $name is back online."
    done
else
    log "System is in desired state. No reconciliation needed."
fi

log "Reconcile completed successfully."
