#!/bin/bash
set -e

echo "[m3tal] Setting up directories..."

mkdir -p /etc/m3tal
mkdir -p /var/lib/m3tal

chmod 755 /etc/m3tal
chmod 755 /var/lib/m3tal

echo "[m3tal] Install complete"
