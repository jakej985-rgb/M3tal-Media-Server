#!/bin/bash

# M3TAL Core Orchestrator (run.sh)
# Logic: Sense -> Think -> Act loop

LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_DIR/core.log"
}

run_agents() {
  log "Loop Start: Sensing..."
  bash agents/monitor.sh
  bash agents/metrics.sh
  
  log "Thinking..."
  bash agents/anomaly-agent.sh
  bash agents/decision-engine.sh
  
  log "Acting..."
  bash agents/reconcile.sh
  bash agents/registry.sh
  
  log "Loop End."
}

log "🚀 Starting m3tal-core orchestrator..."

while true; do
  run_agents
  sleep 10
done
