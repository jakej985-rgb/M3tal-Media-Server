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
    systemctl daemon-reload >/dev/null 2>&1 || true
fi

# 2. Remove UX symlink if it points to our stack
if [ -L /docker ]; then
    LINK_TARGET=$(readlink /docker)
    if [ "$LINK_TARGET" == "/opt/m3tal/stack" ]; then
        echo "[m3tal] Removing /docker symlink"
        rm /docker
    fi
fi

echo "[m3tal] Removed."
echo "[m3tal] NOTE: Configuration (/etc/m3tal) and Data (/var/lib/m3tal) have been preserved."
echo "[m3tal] NOTE: Stack files (/opt/m3tal/stack) have also been preserved."