# M3TAL Core Architecture (Go-Native)

## 🧠 Sense-Think-Act Loop
The system is built on a 6-pillar architecture implemented as high-performance Go agents.

### 1. Monitor (`monitorAgent`)
**Role**: Senses container state.
- Tracks Docker container status (running, exited, dead).
- Gathers per-container CPU, Memory, and I/O metrics.
- Discovers "managed" containers via labels (`m3tal.managed`).

### 2. Metrics (`metricsAgent`)
**Role**: Senses hardware and host state.
- Collects system-wide CPU and Memory usage.
- Discovers hardware temperatures (hwmon/thermal zones).
- Polls GPU metrics (Load, VRAM, Temp) via `radeontop` or sysfs.
- Monitors storage health (Usage, SMART temperatures).

### 3. Anomaly (`anomalyAgent`)
**Role**: Identifies abnormal patterns.
- Watches Docker event streams for sudden crashes or restart loops.
- Scans container logs for critical errors, panics, or fatal exceptions.
- Surfaces "Anomalies" to the state machine.

### 4. Decision (`decisionAgent`)
**Role**: The Brain.
- Analyzes the current state (Health Score).
- Evaluates anomalies and determines corrective actions.
- Enforces cooldowns to prevent flapping.
- Generates "Decisions" for the Reconciler.

### 5. Reconcile (`reconcileAgent`)
**Role**: The Hands.
- Executes approved decisions.
- Restarts failed services.
- Triggers alerts via Notify.
- (Future) Performs auto-scaling actions.

### 6. Registry (`registryAgent`)
**Role**: Memory and Persistence.
- Synchronizes the in-memory state to atomic JSON files.
- Ensures the Dashboard has a real-time view of the system.
- Maintains historical metrics for trend analysis.

---

## 🛠️ Implementation Specs

### Directory Structure
```
source/go-backend/
├── main.go        # Unified agent entrypoint
├── vendor/        # Pinned dependencies
└── logs/          # Isolated agent logs
```

### Persistence
All state is persisted to `./state/*.json` to ensure the Control Plane can restart without losing context.

### Observability
Each pillar logs to its own file:
- `monitor.log`
- `metrics.log`
- `anomaly.log`
- `decision.log`
- `reconcile.log`
- `registry.log`
