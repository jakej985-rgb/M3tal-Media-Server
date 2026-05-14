#!/bin/bash
set -e

echo "[m3tal] Installing M3TAL Core..."

# Create required directories
mkdir -p /etc/m3tal
mkdir -p /var/lib/m3tal/state
mkdir -p /var/log/m3tal

# Create default env file if not present
if [ ! -f /etc/m3tal/env ]; then
    cat > /etc/m3tal/env <<'EOF'
# M3TAL Core Environment Variables
# Override defaults from /etc/m3tal/config.yaml here
M3TAL_LOG_LEVEL=info
M3TAL_DATA_DIR=/var/lib/m3tal
M3TAL_API_PORT=8080
EOF
fi

# Reload systemd and enable services (best-effort — may not be in container)
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload >/dev/null 2>&1 || true
    systemctl enable m3tal.service >/dev/null 2>&1 || true
    systemctl enable m3tal-api.service >/dev/null 2>&1 || true
    echo "[m3tal] Systemd services enabled: m3tal, m3tal-api"
    echo "[m3tal] Start services with: sudo systemctl start m3tal m3tal-api"
fi

echo "[m3tal] ✅ Installation complete."
echo "[m3tal]   CLI:    m3tal"
echo "[m3tal]   API:    m3tal-api"
echo "[m3tal]   Config: /etc/m3tal/config.yaml"
echo "[m3tal]   Stack:  /usr/share/m3tal/stack/"
echo ""
echo "[m3tal] Optional — Dashboard (Docker):"
echo "  docker compose -f /usr/share/m3tal/stack/m3tal-compose.yml up dash"