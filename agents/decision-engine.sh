#!/bin/bash

# M3TAL Decision Engine
# Responsibility: Decide corrective actions

LOG_FILE="logs/decision.log"
STATE_DIR="state"
mkdir -p logs "$STATE_DIR"

log() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a "$LOG_FILE"
}

log "Evaluating system state for decisions..."

if [ -f "state/anomalies.json" ] && grep -q "container_crash" state/anomalies.json; then
    log "PLAN: Triggering self-healing for crashed containers."
    echo "{\"action\": \"reconcile_containers\", \"reason\": \"anomalies_detected\", \"timestamp\": $(date +%s)}" > state/decisions.json
else
    log "No actions required."
    echo "[]" > state/decisions.json
fi

log "Decision phase completed."
