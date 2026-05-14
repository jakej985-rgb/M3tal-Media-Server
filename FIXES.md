# Audit Report: M3TAL Core Documentation
**Auditor:** DocCritic, Senior DevOps Auditor  
**Date:** 2023-10-27  
**Status:** **FAILED**

---

### **Verdict**
**Non-Deployable.** The documentation suffers from "developer myopia." It assumes the user has existing knowledge of the `m3tal` binary's internal dependencies, storage requirements, and environment variable schema. A new user following these instructions will encounter multiple runtime failures immediately after the `up` command.

---

### **Issue List**

#### **BLOCKER**
1.  **Missing `.env` schema validation:** The documentation tells me to `cp template.env .env`, but there is zero documentation on *what* keys are mandatory, especially regarding `BASE_STORAGE_PATH` or external API keys.
2.  **Missing `/mnt` initialization:** The requirement "`/mnt` Directory Writable" is a manual step. If a user on macOS or a restricted Linux distro runs `mkdir /mnt`, they will face permission issues or file system conflicts.
3.  **Port Conflicts/Permissions:** The documentation mandates ports 80 and 443. On Linux, these are privileged ports. The documentation fails to mention that `sudo` or specific `setcap` configurations may be required to bind the Traefik ingress.

#### **WARNING**
4.  **Implicit Build Tooling:** `./build.sh` is invoked, but we don't know what it does. Does it fetch Go dependencies? Does it install Docker SDKs? Does it require `go mod download` first?
5.  **Environment Variable Persistence:** The CLI `config set` is mentioned, but where is it writing to? Does the user need to manually edit `.env` or does the CLI handle it? This is ambiguous.
6.  **"Legacy" Dashboard confusion:** You state the Python dashboard is being phased out, yet you provide instructions for it. Which one should a new user start with?

#### **SUGGESTION**
7.  **DNS Mapping:** Expecting a new user to manually edit `/etc/hosts` is brittle. Consider suggesting a tool like `dnsmasq` or providing a warning about `sudo` requirements for file edits.
8.  **Orchestrator feedback:** The `init` command is mysterious. What state is it creating? Does it check for the existence of `BASE_STORAGE_PATH` before running?

---

### **Suggested Fixes**

*   **For `.env`:** Include a table in the README detailing: `KEY`, `REQUIRED (Y/N)`, `DESCRIPTION`, and `EXAMPLE`. 
*   **For `/mnt`:** Update the instructions to allow configuration of the base path via an environment variable that *doesn't* require writing to root-owned directories. If you *must* use `/mnt`, provide a script: `sudo mkdir -p /mnt/m3tal && sudo chown $USER:$USER /mnt/m3tal`.
*   **For Port 80/443:** Clearly state: *"Note: Binding to ports 80/443 requires root privileges or explicit capability assignment. If running without root, change the ports in your `.env` file to 8080/8443."*
*   **For `build.sh`:** Replace `./build.sh` with a standard `go build -o m3tal ./cmd/main.go` or document what the script handles to ensure the user isn't blindly executing unknown scripts.
*   **For Service Selection:** Clearly label the `m3tal-godash` as the recommended path for new users and explicitly mark the Python dashboard as "Deprecated/Legacy."

---

**Auditor Note:** *I cannot recommend this project for production or even lab deployment until the bootstrap process is automated and the environment constraints are explicitly documented. The current README assumes the user is the original developer.*