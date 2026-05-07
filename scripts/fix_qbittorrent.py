import os
import re
from pathlib import Path

# --- Configuration ---
# This script targets the qBittorrent configuration file to resolve "Access Denied" 
# or connectivity issues when running behind Traefik/VPN.

def fix_qbittorrent_config():
    # Attempt to locate DATA_DIR from .env
    repo_root = Path(__file__).resolve().parent.parent
    env_file = repo_root / ".env"
    data_dir = None

    if env_file.exists():
        with open(env_file, "r") as f:
            for line in f:
                if line.startswith("DATA_DIR="):
                    data_dir = line.split("=", 1)[1].strip()
                    break

    if not data_dir:
        print("[X] Could not find DATA_DIR in .env. Falling back to default.")
        data_dir = str(repo_root / "media")

    # Normalize path (handle Windows/Linux differences)
    data_path = Path(data_dir.replace("\\", "/"))
    config_file = data_path / "config" / "qbittorrent" / "qBittorrent" / "qBittorrent.conf"

    if not config_file.exists():
        print(f"[X] Config file not found at: {config_file}")
        print("    Ensure qBittorrent has started at least once.")
        return

    print(f"[*] Found config at: {config_file}")
    
    with open(config_file, "r") as f:
        content = f.read()

    # Settings to enforce
    updates = {
        r"WebUI\\HostHeaderValidation=.*": "WebUI\\HostHeaderValidation=false",
        r"WebUI\\CSRFProtection=.*": "WebUI\\CSRFProtection=false",
        r"WebUI\\Address=.*": "WebUI\\Address=0.0.0.0",
        r"WebUI\\Port=.*": "WebUI\\Port=8090",
        r"WebUI\\LocalHostAuthentication=.*": "WebUI\\LocalHostAuthentication=false"
    }

    modified = content
    for pattern, replacement in updates.items():
        if re.search(pattern, modified):
            modified = re.sub(pattern, replacement, modified)
        else:
            # If the setting isn't there, add it to the [WebUI] section
            if "[WebUI]" in modified:
                modified = modified.replace("[WebUI]", f"[WebUI]\n{replacement}")
            else:
                modified += f"\n[WebUI]\n{replacement}\n"

    if modified != content:
        with open(config_file, "w") as f:
            f.write(modified)
        print("[+] qBittorrent configuration updated successfully.")
        print("[!] Please restart the qbittorrent container for changes to take effect.")
        print("    Command: docker restart qbittorrent")
    else:
        print("[i] Configuration already contains the recommended settings.")

if __name__ == "__main__":
    fix_qbittorrent_config()
