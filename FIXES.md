### **DocCritic Audit Report: M3TAL Core v1.7**

**Audit Status:** FAILED
**Verdict:** **BLOCKER.** The documentation assumes a level of prior domain knowledge and environmental state that does not exist for a fresh user. The "Quick Start" is currently a path to a broken deployment.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing `.env` Validation/Configuration:** The user is told to `cp template.env .env`, but there is zero instruction on what variables *must* be changed (specifically `BASE_STORAGE_PATH`). If the user runs `./m3tal up` without configuring the storage path, the container will likely crash or mount nothing.
2.  **Implicit Host Requirement:** The instruction "Your host data at `BASE_STORAGE_PATH` is always mounted to `/mnt`" implies that the directory must exist on the host before running `./m3tal up`. If it doesn't exist, Docker will create a root-owned directory, causing permission issues.
3.  **Missing `m3tal` Binary Context:** The instructions jump from `build.sh` to `./m3tal init`. It is not explicitly stated that `build.sh` produces the `m3tal` binary, nor what the user should do if that binary fails to generate (e.g., Go dependency errors).

#### **WARNING**
4.  **Incomplete Networking Instructions:** While the `/etc/hosts` file is mentioned, the Traefik configuration dependency is not. If Traefik requires a specific network name or certificate setup to bind to those hosts, the user is left guessing.
5.  **Ambiguous Dashboard Status:** You mention the legacy Dashboard is being phased out, yet you provide it as a main access point in the routing table. A new user will be confused about whether they should use `m3tal-godash` or the legacy dashboard.

#### **SUGGESTION**
6.  **CLI Help Accessibility:** The documentation should explicitly mention `./m3tal --help` as a way to verify the binary is functional after the build step.
7.  **Dependency Verification:** Add a step to check for the Docker Socket permissions. `docker` group membership is mentioned, but a new user often forgets to log out/in to apply group changes.

---

### **Suggested Fixes**

#### **1. Improve `Prerequisites` section:**
Add a "Host Preparation" step:
```bash
# Verify Docker socket permissions
groups | grep docker || echo "WARNING: User not in docker group. Please add user and restart session."

# Prepare Storage
mkdir -p /path/to/your/media
# Ensure your .env BASE_STORAGE_PATH points to this directory
```

#### **2. Update `Quick Start` to be "Configuration-First":**
```bash
# 1. Setup
git clone ...
cp template.env .env
# IMPORTANT: Edit .env and set BASE_STORAGE_PATH to your absolute media directory.
nano .env

# 2. Build & Verify
chmod +x build.sh
./build.sh
./m3tal --version # Verify binary is executable
```

#### **3. Clarify the Storage/Pathing behavior:**
Add a warning to the `Troubleshooting` section:
> **Crucial:** M3TAL expects `BASE_STORAGE_PATH` to be an absolute path. If the directory does not exist, the orchestrator will fail to bind mount, resulting in an empty `/mnt` inside your containers. 

#### **4. Address the Dashboard confusion:**
Add a note under "Service Routing":
> **Note:** The legacy Dashboard (`m3tal.localhost`) is included for v1.7. If you are deploying for the first time, we recommend checking the [m3tal-godash](https://github.com/jakej985-rgb/m3tal-godash) repository for the modern replacement.

---
**Auditor Signature:** *DocCritic*
*M3TAL Platform Audit Division*