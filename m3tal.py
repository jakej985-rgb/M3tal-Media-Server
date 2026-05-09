import sys
import os
import subprocess
import argparse
import re
from pathlib import Path

# Batch 16 Hardening: Force UTF-8 for Windows console resilience
if sys.stdout.encoding.lower() != 'utf-8':
    try:
        sys.stdout.reconfigure(encoding='utf-8')
    except (AttributeError, Exception):
        pass

# M3TAL Unified CLI (v2.4 Production-Grade)
# Responsibility: Centralized entrypoint for all M3TAL orchestration.

# --- Environment Variable Bootstrap -------------------------------------------
def _bootstrap_env():
    # Attempt to find root to locate .env
    p = Path(__file__).resolve()
    root = None
    for parent in [p] + list(p.parents):
        if (parent / ".env").exists() and (parent / "docker").exists():
            root = parent
            break
    if not root:
        return
    
    with open(root / ".env", "r", encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("#"):
                continue
            if line.startswith("export "):
                line = line[len("export "):].strip()
            if "=" in line:
                k, v = line.split("=", 1)
                k, v = k.strip(), v.strip()
                v = re.sub(r'\s+#.*$', '', v).strip()
                if (v.startswith('"') and v.endswith('"')) or (v.startswith("'") and v.endswith("'")):
                    v = v[1:-1]
                os.environ[k] = v

    # Audit Fix: Enforce Project Root
    os.chdir(root)
    return root

ROOT = _bootstrap_env()

# --- Execution Helpers --------------------------------------------------------
def get_compose_files():
    """Dynamically discover all docker-compose files in priority order."""
    docker_dir = ROOT / "docker"
    compose_files = []
    
    # Add m3tal-compose first (since it has the proxy and backend)
    cp_compose = docker_dir / "m3tal-compose.yml"
    if cp_compose.exists():
        compose_files.append(cp_compose)
        
    # Then network and routing
    network_compose = docker_dir / "network-compose.yml"
    if network_compose.exists():
        compose_files.append(network_compose)
        
    routing_compose = docker_dir / "routing-compose.yml"
    if routing_compose.exists():
        compose_files.append(routing_compose)
        
    # Then the rest dynamically
    if docker_dir.exists():
        for file in docker_dir.glob("*-compose.yml"):
            if file not in compose_files:
                compose_files.append(file)
                
    # Also discover any docker-compose.yml in root of docker/ if any
    root_compose = docker_dir / "docker-compose.yml"
    if root_compose.exists() and root_compose not in compose_files:
        compose_files.append(root_compose)

    return compose_files

# --- Command Handlers ---------------------------------------------------------
def cmd_obsolete():
    print("[!] This command is obsolete and has been removed in the Go-native M3TAL Control Plane.")
    return 1

def cmd_build(args):
    """Enforces a clean rebuild of the control plane containers."""
    print("\n[BUILD] Triggering no-cache build of M3TAL Control Plane...")
    compose_file = ROOT / "docker" / "m3tal-compose.yml"
    cmd = ["docker", "compose", "-f", str(compose_file), "build", "--no-cache"]
    try:
        subprocess.run(cmd, cwd=str(ROOT / "docker"), check=True)
        print("[INIT] Build successful. Containers are up to date.")
        return 0
    except Exception as e:
        print(f"[X] Build failed: {e}")
        return 1

def cmd_init(args):
    """Initializes the M3TAL environment."""
    compose_files = get_compose_files()
    if not compose_files:
        print("[X] No compose files found!")
        return 1
        
    for cf in compose_files:
        print(f"\n[INIT] Starting stack: {cf.name}...")
        cmd = ["docker", "compose", "-f", str(cf), "up", "-d"]
        try:
            subprocess.run(cmd, cwd=str(cf.parent), check=True)
        except Exception as e:
            print(f"[X] Failed to start {cf.name}: {e}")
            return 1
            
    print("\n[INIT] M3TAL environment initialization complete.")
    return 0

def cmd_shutdown(args):
    """Executes the Global Blackout or selective shutdown."""
    compose_files = get_compose_files()
    if not compose_files:
        print("[X] No compose files found!")
        return 1
        
    # Reverse the order for shutdown
    compose_files_rev = compose_files.copy()
    compose_files_rev.reverse()
    
    for cf in compose_files_rev:
        print(f"\n[SHUTDOWN] Stopping stack: {cf.name}...")
        cmd = ["docker", "compose", "-f", str(cf), "down"]
        try:
            subprocess.run(cmd, cwd=str(cf.parent), check=True)
        except Exception as e:
            print(f"[X] Failed to stop {cf.name}: {e}")
            
    print("\n[SHUTDOWN] M3TAL environment shutdown complete.")
    return 0

def cmd_ps(args):
    """Shows the status of all M3TAL containers."""
    compose_files = get_compose_files()
    if not compose_files:
        print("[X] No compose files found!")
        return 1
    for cf in compose_files:
        print(f"\n[STATUS] Stack: {cf.name}")
        subprocess.run(["docker", "compose", "-f", str(cf), "ps"], cwd=str(cf.parent))
    return 0

def cmd_restart(args):
    """Restarts the M3TAL environment."""
    print("\n[RESTART] Triggering full system restart...")
    cmd_shutdown(args)
    return cmd_init(args)

# --- CLI Structure ------------------------------------------------------------
def main():
    if not ROOT:
        print("[X] FATAL: Could not locate M3TAL repository root.")
        print("   Please run this from within the M3TAL project folder.")
        sys.exit(1)

    parser = argparse.ArgumentParser(
        prog="m3tal",
        description="M3TAL Control Plane - Production Tooling CLI"
    )
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")

    # up / start
    for cmd_name in ["up", "start"]:
        p = subparsers.add_parser(cmd_name, help="Initialize and start the M3TAL environment")
        p.add_argument("--repair", help="Legacy argument (ignored)")

    # down / stop
    for cmd_name in ["down", "stop"]:
        p = subparsers.add_parser(cmd_name, help="Safely stop all M3TAL stacks and services")
        p.add_argument("stacks", nargs="*", help="Legacy argument (ignored)")

    # restart
    subparsers.add_parser("restart", help="Restart the M3TAL environment")

    # ps / ls / status
    for cmd_name in ["ps", "ls", "status"]:
        subparsers.add_parser(cmd_name, help="Show status of M3TAL containers")

    # build
    subparsers.add_parser("build", help="Enforce no-cache rebuild of M3TAL Docker images")

    # obsolete commands to maintain CLI interface without crashing
    for cmd in ["logs", "env", "audit", "traefik", "test", "run", "heal", "config", "dashpass"]:
        p = subparsers.add_parser(cmd, help="[Obsolete in Go-native version]")
        p.add_argument("args", nargs=argparse.REMAINDER)

    # If no args, show help
    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(0)

    args = parser.parse_args()

    # Context Guard: Ensure we have a valid environment before any command
    if not (ROOT / ".env").exists():
        print("[X] FATAL: Missing .env file at repository root.")
        print(f"   Searching in: {ROOT}")
        sys.exit(1)

    if args.command in ["up", "start"]:
        sys.exit(cmd_init(args))
    elif args.command in ["down", "stop"]:
        sys.exit(cmd_shutdown(args))
    elif args.command == "restart":
        sys.exit(cmd_restart(args))
    elif args.command in ["ps", "ls", "status"]:
        sys.exit(cmd_ps(args))
    elif args.command == "build":
        sys.exit(cmd_build(args))
    else:
        sys.exit(cmd_obsolete())

if __name__ == "__main__":
    main()
