import docker
import json
import os
import re
import time
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s [%(levelname)s] %(message)s')
logger = logging.getLogger("network-scout")

STATE_DIR = os.getenv('STATE_DIR', '/docker/control-plane/state')
NETWORK_JSON = os.path.join(STATE_DIR, 'network.json')

def scout_routes():
    """Scan Docker containers for Traefik routing labels and build a link list."""
    try:
        client = docker.from_env()
        containers = client.containers.list()
        links = []
        seen_hosts = set()
        
        # Blacklist of internal/system names
        blacklist = ['dashboard', 'api', 'traefik', 'docker-proxy', 'glances', 'dozzle', 'portainer']
        
        for container in containers:
            labels = container.labels
            # Look for Traefik host rules
            for key, value in labels.items():
                if 'traefik.http.routers.' in key and '.rule' in key:
                    # Extract Host(...) from the rule
                    match = re.search(r'Host\(`([^`]+)`\)', value)
                    if match:
                        host = match.group(1)
                        if host in seen_hosts:
                            continue
                            
                        # Determine a readable name
                        name = host.split('.')[0].capitalize()
                        
                        # Skip blacklisted items
                        if any(b.lower() in name.lower() for b in blacklist):
                            continue
                            
                        links.append({
                            "name": name,
                            "url": f"https://{host}",
                            "status": "enabled",
                            "container": container.name
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
