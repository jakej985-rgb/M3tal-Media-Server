#!/bin/bash
set -e

echo "[m3tal] Removing M3TAL Core..."

# 1. Stop and disable services
if command -v systemctl >/dev/null 2>&1; then
    if [ -f /lib/systemd/system/m3tal-api.service ]; then
        echo "[m3tal] Stopping M3TAL API service"
        systemctl stop m3tal-api.service >/dev/null 2>&1 || true
        systemctl disable m3tal-api.service >/dev/null 2>&1 || true
    fi
    if [ -f /lib/systemd/system/m3tal.service ]; then
        echo "[m3tal] Stopping M3TAL Agents service"
        systemctl stop m3tal.service >/dev/null 2>&1 || true
        systemctl disable m3tal.service >/dev/null 2>&1 || true
    fi
    systemctl daemon-reload >/dev/null 2>&1 || true
fi

# 4. Clean up /docker symlink
if [ -L /docker ]; then
    echo "[m3tal] Removing /docker symlink"
    rm /docker
fi

# 5. Handle Purge (remove data and config)
if [ "$1" = "purge" ]; then
    echo "[m3tal] Purging all M3TAL data and configuration..."
    rm -rf /opt/m3tal
    rm -rf /etc/m3tal
    rm -rf /var/lib/m3tal
    echo "[m3tal] Purge complete."
else
    echo "[m3tal] Removed."
    echo "[m3tal] NOTE: Configuration (/etc/m3tal) and Data (/var/lib/m3tal) have been preserved."
    echo "[m3tal] NOTE: Stack files (/opt/m3tal/stack) have also been preserved."
fi