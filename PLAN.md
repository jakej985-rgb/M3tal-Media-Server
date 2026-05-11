# m3tal-core Plan (Control Plane)

## Purpose
Central brain of the system.

Handles:
- Agents
- Automation
- Self-healing
- System state

---

## Goals
- Single entrypoint (`run.sh`)
- Stable agent execution loop
- Centralized logging + state
- No Docker definitions here

---

## Structure

m3tal-core/
├── m3tal.py
├── agents/
├── logs/
├── state/
├── config/
├── internal/
└── run.sh

---

## Tasks

### 1. Normalize agents
- Move all agents into `/agents`
- Ensure all scripts are executable

chmod +x agents/*.sh

---

### 2. Fix run loop

Create `run.sh`:

#!/bin/bash

LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_DIR/core.log"
}

run_agents() {
  bash agents/monitor.sh
  bash agents/metrics.sh
  bash agents/anomaly-agent.sh
  bash agents/decision-engine.sh
  bash agents/reconcile.sh
  bash agents/registry.sh
}

log "Starting m3tal-core..."

while true; do
  run_agents
  sleep 10
done

---

### 3. Standardize paths

ALL agents must use:

/mnt

---

### 4. Permissions

sudo chown -R 1000:1000 /mnt
sudo chmod -R 775 /mnt

---

### 5. Logging

Each agent logs to:

/logs/<agent>.log

---

### 6. State system (required for API later)

Create:

state/system.json

Example:

{
  "status": "healthy",
  "containers": {},
  "last_check": ""
}

---

## Done When

- `./run.sh` runs continuously
- No agent crashes
- Logs populate
- State file updates
