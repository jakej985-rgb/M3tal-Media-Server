# 🛠️ M3TAL CLI Guide (m3tal.py)

The `m3tal.py` script is the unified entry point for all system orchestration. It simplifies complex Docker Compose commands and ensures cross-platform path compatibility.

---

## 🚀 Lifecycle Management
These commands control the operational state of the entire M3TAL environment.

| Command | Aliases | Description |
| :--- | :--- | :--- |
| **`up`** | `start` | Initializes and starts all Docker stacks in priority order. |
| **`down`** | `stop` | Safely stops and removes all M3TAL containers and networks. Supports `--remove-orphans`. |
| **`restart`** | — | Performs a full `down` followed by an `up` cycle. |
| **`pull`** | — | Pulls the latest images for all stacks from GHCR (Registry). |
| **`build`** | — | Enforces a `--no-cache` rebuild of all local Docker images. |
| **`dashpass`** | — | Resets the password for a dashboard user (e.g., `admin`). |

---

## 📊 Observability
Use these commands to check the health and status of the system.

| Command | Aliases | Description |
| :--- | :--- | :--- |
| **`ps`** | `ls`, `status` | Lists all active containers across all M3TAL stacks with status and ports. |

---

## 🛡️ Obsolete Commands
These commands were part of the legacy Python infrastructure. They are retained in the CLI to prevent script breakage but will return a notice informing you of their replacement by the **Go-native Control Plane**.

*   `logs`: Replaced by native Docker log forwarding.
*   `env`: Replaced by `scripts/view_env.py`.
*   `audit`: Replaced by native Go auditing.
*   `heal`: Replaced by the Go **Healer** service.
*   `scan`: Replaced by the Go **Intelligence Feed**.
*   `config`: Handled via `.env` and `configure_env.py`.

---

## 💡 Pro-Tip: Global Alias
For the best experience, add a global alias to your shell:

**Linux (bash/zsh):**
```bash
echo 'alias m3tal="python3 m3tal.py"' >> ~/.bashrc
source ~/.bashrc
```

**Windows (PowerShell Profile):**
```powershell
function m3tal { python m3tal.py $args }
```

Now you can simply run `m3tal up` or `m3tal status` from anywhere in the repository! 🚀
