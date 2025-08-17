#!/bin/bash
set -e

BUN_VERSION="1.2.20"
echo "Installing Bun..."

export BUN_INSTALL=/usr/local

curl -fsSL https://bun.com/install | bash -s "bun-v${BUN_VERSION}"
