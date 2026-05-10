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
    # Attempt to find root by looking for the 'docker' folder
    p = Path(__file__).resolve()
    root = None
    for parent in [p] + list(p.parents):
        if (parent / "docker").exists() and (parent / "m3tal.py").exists():
            root = parent
            break
    if not root:
        return
    
    # Load .env only if it exists (might be missing in CI)
    env_path = root / ".env"
    if env_path.exists():
        with open(env_path, "r", encoding="utf-8") as f:
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
    
    # 1. Network setup MUST be first
    network_compose = docker_dir / "network-compose.yml"
    if network_compose.exists():
        compose_files.append(network_compose)

    # 2. Routing/Traefik
    routing_compose = docker_dir / "routing-compose.yml"
    if routing_compose.exists():
        compose_files.append(routing_compose)
        
    # 3. M3TAL Core (Dashboard, Go-Backend)
    cp_compose = docker_dir / "m3tal-compose.yml"
    if cp_compose.exists():
        compose_files.append(cp_compose)
        
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

def get_project_name(cf_path):
    """Derives a clean project name from the compose filename."""
    name = cf_path.stem.replace("-compose", "").replace(".yml", "")
    if name == "docker": # Handle docker-compose.yml
        return "m3tal-core"
    return f"m3tal-{name}"

# --- Command Handlers ---------------------------------------------------------
def cmd_obsolete():
    print("[!] This command is obsolete and has been removed in the Go-native M3TAL Control Plane.")
    return 1

def cmd_build(args):
    """Enforces a clean rebuild of the control plane containers."""
    print("\n[BUILD] Triggering no-cache build of M3TAL Control Plane...")
    compose_file = ROOT / "docker" / "m3tal-compose.yml"
    project_name = get_project_name(compose_file)
    env_file = ROOT / ".env"
    cmd = ["docker", "compose", "-p", project_name, "--env-file", str(env_file), "-f", str(compose_file), "build", "--no-cache"]
    try:
        subprocess.run(cmd, cwd=str(ROOT / "docker"), check=True)
        print("[INIT] Build successful. Containers are up to date.")
        return 0
    except Exception as e:
        print(f"[X] Build failed: {e}")
        return 1

def _ensure_network():
    """Ensures the M3TAL 'proxy' network exists."""
    print("[INIT] Ensuring 'proxy' network exists...")
    try:
        subprocess.run(["docker", "network", "inspect", "proxy"], capture_output=True, check=True)
    except subprocess.CalledProcessError:
        print("[INIT] Creating 'proxy' network...")
        try:
            subprocess.run(["docker", "network", "create", "proxy"], check=True)
        except Exception as e:
            print(f"[X] Failed to create network: {e}")
            return False
    return True

def cmd_init(args):
    """Initializes the M3TAL environment."""
    if not _ensure_network():
        return 1
        
    compose_files = get_compose_files()
    if not compose_files:
        print("[X] No compose files found!")
        return 1
        
    target = getattr(args, "stack", None)
    
    for cf in compose_files:
        if target and target.lower() not in cf.name.lower():
            continue
            
        print(f"\n[INIT] Starting stack: {cf.name}...")
        project_name = get_project_name(cf)
        env_file = ROOT / ".env"
        cmd = ["docker", "compose", "-p", project_name, "--env-file", str(env_file), "-f", str(cf), "up", "-d"]
        
        if getattr(args, "recreate", False):
            cmd.append("--force-recreate")
        if getattr(args, "remove_orphans", False):
            cmd.append("--remove-orphans")
            
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
    
    target = getattr(args, "stack", None)
    
    for cf in compose_files_rev:
        if target and target.lower() not in cf.name.lower():
            continue
            
        print(f"\n[SHUTDOWN] Stopping stack: {cf.name}...")
        project_name = get_project_name(cf)
        env_file = ROOT / ".env"
        cmd = ["docker", "compose", "-p", project_name, "--env-file", str(env_file), "-f", str(cf), "down"]
        
        if getattr(args, "remove_orphans", False):
            cmd.append("--remove-orphans")
            
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
    target = getattr(args, "stack", None)
    
    for cf in compose_files:
        if target and target.lower() not in cf.name.lower():
            continue
            
        print(f"\n[STATUS] Stack: {cf.name}")
        project_name = get_project_name(cf)
        env_file = ROOT / ".env"
        subprocess.run(["docker", "compose", "-p", project_name, "--env-file", str(env_file), "-f", str(cf), "ps"], cwd=str(cf.parent))
    return 0

def cmd_restart(args):
    """Restarts the M3TAL environment using recreation (up -r)."""
    print("\n[RESTART] Triggering system recreation (up --force-recreate)...")
    args.recreate = True
    # Default to remove orphans on restart for cleanliness
    if not hasattr(args, "remove_orphans"):
        args.remove_orphans = True
    return cmd_init(args)

def cmd_pull(args):
    """Pulls the latest images for the M3TAL stacks."""
    compose_files = get_compose_files()
    if not compose_files:
        print("[X] No compose files found!")
        return 1
        
    # If a specific stack was requested, filter for it
    target = args.stack if hasattr(args, "stack") else None
    
    for cf in compose_files:
        if target and target.lower() not in cf.name.lower():
            continue
            
        print(f"\n[PULL] Refreshing images for: {cf.name}...")
        project_name = get_project_name(cf)
        env_file = ROOT / ".env"
        cmd = ["docker", "compose", "-p", project_name, "--env-file", str(env_file), "-f", str(cf), "pull"]
        try:
            subprocess.run(cmd, cwd=str(cf.parent), check=True)
        except Exception as e:
            print(f"[X] Failed to pull {cf.name}: {e}")
            
    print("\n[PULL] Image synchronization complete.")
    return 0

def cmd_dashpass(args):
    """Manages dashboard users and passwords interactively."""
    import json
    try:
        import bcrypt
    except ImportError:
        print("[X] Error: 'bcrypt' is required for password hashing.")
        print("    Install it with: pip install bcrypt")
        return 1
        
    users_file = ROOT / "docker" / "state" / "users.json"
    
    # Audit Fix: Ensure parent directory exists
    users_file.parent.mkdir(parents=True, exist_ok=True)
    
    # Load existing users
    try:
        with open(users_file, "r") as f:
            users = json.load(f)
    except Exception:
        users = []
        
    # Automatic mode if arguments are provided via CLI
    if args.password:
        username = args.username or "admin"
        password = args.password
        print(f"\n[AUTH] Non-interactive reset for '{username}'...")
    else:
        # Interactive Menu
        print("\n--- M3TAL User Management ---")
        print("1.) Reset 'admin' password")
        print("2.) Manage existing users")
        print("3.) Add new user")
        print("q.) Quit")
        
        choice = input("\nSelect an option: ").strip().lower()
        
        username = None
        if choice == "1":
            username = "admin"
        elif choice == "2":
            if not users:
                print("[!] No users found.")
                return 0
            print("\nExisting Users:")
            for i, u in enumerate(users):
                print(f"  {i+1}.) {u['username']} [{u.get('role', 'viewer')}]")
            idx = input("\nSelect user index: ").strip()
            try:
                username = users[int(idx)-1]['username']
            except (ValueError, IndexError):
                print("[X] Invalid selection.")
                return 1
        elif choice == "3":
            username = input("Enter new username: ").strip()
            if not username:
                return 0
        elif choice == "q":
            return 0
        else:
            print("[X] Invalid choice.")
            return 1

        import getpass
        password = getpass.getpass(f"Enter password for '{username}': ")
        confirm = getpass.getpass("Confirm password: ")
        if password != confirm:
            print("[X] Passwords do not match!")
            return 1
            
    # Hash and Save
    pwd_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')
    
    # Update record
    role = "admin" if username == "admin" else "viewer"
    # Find if user exists to preserve role if not admin
    existing = next((u for u in users if u.get("username") == username), None)
    if existing and username != "admin":
        role = existing.get("role", "viewer")
        
    users = [u for u in users if u.get("username") != username]
    users.append({
        "username": username,
        "token_hash": pwd_hash,
        "role": role
    })
    
    with open(users_file, "w") as f:
        json.dump(users, f, indent=2)
        
    print(f"\n[AUTH] Success! '{username}' has been updated.")
    print("[AUTH] Changes are live (no dashboard restart required).")
    return 0

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
        p.add_argument("stack", nargs="?", help="Specific stack to up (e.g. m3tal, media)")
        p.add_argument("-r", "--recreate", action="store_true", help="Force recreate containers")
        p.add_argument("--remove-orphans", action="store_true", help="Remove containers for services not defined in the Compose file")

    # down / stop
    for cmd_name in ["down", "stop"]:
        p = subparsers.add_parser(cmd_name, help="Safely stop all M3TAL stacks and services")
        p.add_argument("stack", nargs="?", help="Specific stack to stop (e.g. m3tal, media)")
        p.add_argument("--remove-orphans", action="store_true", help="Remove containers for services not defined in the Compose file")

    # restart
    p_restart = subparsers.add_parser("restart", help="Restart the M3TAL environment")
    p_restart.add_argument("stack", nargs="?", help="Specific stack to restart (e.g. m3tal, media)")
    p_restart.add_argument("--remove-orphans", action="store_true", help="Remove containers for services not defined in the Compose file during shutdown")

    # ps / ls / status
    for cmd_name in ["ps", "ls", "status"]:
        p = subparsers.add_parser(cmd_name, help="Show status of M3TAL containers")
        p.add_argument("stack", nargs="?", help="Specific stack to show (e.g. m3tal, media)")

    # build
    subparsers.add_parser("build", help="Enforce no-cache rebuild of M3TAL Docker images")

    # pull [stack]
    p_pull = subparsers.add_parser("pull", help="Pull latest images from registry (GHCR)")
    p_pull.add_argument("stack", nargs="?", help="Specific stack to pull (e.g. m3tal, media)")

    # dashpass [username] [--password PWD]
    p_pass = subparsers.add_parser("dashpass", help="Reset dashboard admin password")
    p_pass.add_argument("username", nargs="?", default="admin", help="Username to reset (default: admin)")
    p_pass.add_argument("--password", help="Set password non-interactively")

    # obsolete commands to maintain CLI interface without crashing
    for cmd in ["logs", "env", "audit", "traefik", "test", "run", "heal", "config"]:
        p = subparsers.add_parser(cmd, help="[Obsolete in Go-native version]")
        p.add_argument("args", nargs=argparse.REMAINDER)

    # If no args, show help
    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(0)

    args = parser.parse_args()

    # Context Guard: Ensure we have a valid environment for lifecycle commands
    env_missing = not (ROOT / ".env").exists()
    if env_missing and args.command in ["up", "start", "restart", "down", "stop"]:
        print("[X] FATAL: Missing .env file at repository root.")
        print("   This command requires a valid environment configuration.")
        sys.exit(1)
    elif env_missing:
        print("[!] Warning: Missing .env file. Some features may not work as expected.")

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
    elif args.command == "pull":
        sys.exit(cmd_pull(args))
    elif args.command == "dashpass":
        sys.exit(cmd_dashpass(args))
    else:
        sys.exit(cmd_obsolete())

if __name__ == "__main__":
    main()
