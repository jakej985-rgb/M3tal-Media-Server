#!/bin/bash
set -e

echo "[m3tal] Removing M3TAL Core..."

# Stop and disable services (best-effort)
if command -v systemctl >/dev/null 2>&1; then
    systemctl stop m3tal.service >/dev/null 2>&1 || true
    systemctl stop m3tal-api.service >/dev/null 2>&1 || true
    systemctl disable m3tal.service >/dev/null 2>&1 || true
    systemctl disable m3tal-api.service >/dev/null 2>&1 || true
    systemctl daemon-reload >/dev/null 2>&1 || true
    echo "[m3tal] Services stopped and disabled."
fi

echo "[m3tal] Removed (config preserved in /etc/m3tal/, data preserved in /var/lib/m3tal/)"