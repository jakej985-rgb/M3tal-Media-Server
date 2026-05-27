#!/bin/bash
set -e

echo "[m3tal] Running post-installation setup..."

# 1. Create required directories
mkdir -p /etc/m3tal
mkdir -p /var/lib/m3tal
mkdir -p /opt/m3tal/stack
mkdir -p /var/log/m3tal

# Ensure the docker group can manage stacks
chown -R root:docker /opt/m3tal/stack
chmod -R 775 /opt/m3tal/stack

# 2. Config Initialization (do not overwrite existing config)
if [ ! -f /etc/m3tal/.env ]; then
    if [ -f /etc/m3tal/.env.example ]; then
        echo "[m3tal] Initializing default configuration at /etc/m3tal/.env"
        cp /etc/m3tal/.env.example /etc/m3tal/.env
        chmod 600 /etc/m3tal/.env
    else
        echo "⚠️  /etc/m3tal/.env.example not found, skipping config init"
    fi
fi

# 3. Stack Initialization (templates to /opt/m3tal/stack)
if [ -d /usr/share/m3tal/stack ]; then
    echo "[m3tal] Initializing M3TAL Stack files at /opt/m3tal/stack"
    mkdir -p /opt/m3tal/stack/dynamic
    cp -rn /usr/share/m3tal/stack/. /opt/m3tal/stack/
fi

# 4. /docker UX Symlink
if [ ! -e /docker ]; then
    echo "[m3tal] Creating /docker symlink to /opt/m3tal/stack"
    ln -s /opt/m3tal/stack /docker
elif [ ! -L /docker ]; then
    echo "⚠️  /docker exists and is not a symlink, skipping symlink creation"
fi

# 5. Permissions & Group Management
echo "[m3tal] Hardening permissions for sudo-less operation..."
if ! getent group m3tal >/dev/null; then
    groupadd -r m3tal
fi

# Ensure all core directories are owned by the group
chown -R root:m3tal /opt/m3tal /etc/m3tal /var/lib/m3tal 2>/dev/null || true
chmod -R g+rwX /opt/m3tal /etc/m3tal /var/lib/m3tal 2>/dev/null || true

# 5. Systemd Service Integration
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload >/dev/null 2>&1 || true
    if [ -f /lib/systemd/system/m3tal-api.service ]; then
        systemctl enable m3tal-api.service >/dev/null 2>&1 || true
        systemctl restart m3tal-api.service >/dev/null 2>&1 || true
        echo "[m3tal] M3TAL API service enabled and started"
    fi
    if [ -f /lib/systemd/system/m3tal.service ]; then
        systemctl enable m3tal.service >/dev/null 2>&1 || true
        systemctl restart m3tal.service >/dev/null 2>&1 || true
        echo "[m3tal] M3TAL Agents service enabled and started"
    fi
fi

# 6. Desktop & Icon Cache Update
if command -v update-desktop-database >/dev/null 2>&1; then
    echo "[m3tal] Updating desktop database..."
    update-desktop-database -q /usr/share/applications || true
fi

if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    echo "[m3tal] Updating GTK icon cache..."
    gtk-update-icon-cache -f -t /usr/share/icons/hicolor >/dev/null 2>&1 || true
fi

echo "[m3tal] ✅ Installation complete."
echo "[m3tal]   CLI:    /usr/bin/m3tal"
echo "[m3tal]   UX:     /docker"
echo "[m3tal]   Config: /etc/m3tal/.env"
echo ""
echo "[m3tal] Quick Start:"
echo "  1. sudo m3tal init  (to generate .env configuration)"
echo "  2. m3tal            (to open the interactive Control Center)"
echo "  3. m3tal help       (to see all available CLI commands)"