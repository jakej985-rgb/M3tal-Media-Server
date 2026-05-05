import subprocess
import json
import os
import re
import time
import logging

from utils.paths import STATE_DIR
from utils.logger import get_logger

# Configure logging
logger = get_logger("scout")
NETWORK_JSON = os.path.join(STATE_DIR, 'network.json')

def scout_routes():
    """Scan Docker containers for Traefik routing labels using raw Docker CLI."""
    try:
        # Use docker ps -a to get names and labels without needing the 'docker' python library
        cmd = ["docker", "ps", "-a", "--format", "{{.Names}}|{{.Labels}}"]
        res = subprocess.run(cmd, capture_output=True, text=True, timeout=15)
        
        if res.returncode != 0:
            logger.error(f"Docker PS failed: {res.stderr}")
            return
            
        lines = res.stdout.strip().splitlines()
        links = []
        seen_hosts = set()
        
        # Blacklist of internal/system names
        blacklist = ['dashboard', 'api', 'traefik', 'docker-proxy', 'glances', 'dozzle', 'portainer']
        
        for line in lines:
            if '|' not in line: continue
            name, labels = line.split('|', 1)
            
            # Look for Traefik host rules in labels
            # Labels format: label1=val1,label2=val2
            if 'traefik.http.routers.' in labels and '.rule=' in labels:
                # Extract Host(...) from the labels string
                match = re.search(r'Host\(`([^`]+)`\)', labels)
                if match:
                    host = match.group(1)
                    if host in seen_hosts:
                        continue
                        
                    # Determine a readable name
                    readable_name = host.split('.')[0].replace('-', ' ').capitalize()
                    
                    # Skip blacklisted items
                    if any(b.lower() in readable_name.lower() for b in blacklist):
                        continue
                        
                    links.append({
                        "name": readable_name,
                        "url": f"https://{host}",
                        "status": "enabled",
                        "container": name
                    })
                    seen_hosts.add(host)
        
        # Sort alphabetically
        links.sort(key=lambda x: x['name'])
        
        # Atomic write to JSON
        temp_file = NETWORK_JSON + '.tmp'
        with open(temp_file, 'w') as f:
            json.dump(links, f, indent=4)
        os.replace(temp_file, NETWORK_JSON)
        
        logger.info(f"Discovered {len(links)} network routes via Docker labels")
        
    except Exception as e:
        logger.error(f"Network Scout failure: {e}")

if __name__ == "__main__":
    logger.info("M3TAL Network Scout starting...")
    while True:
        scout_routes()
        time.sleep(60)  # Scan every minute
