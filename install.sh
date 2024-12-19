#!/usr/bin/env bash

INSTALL_DIR=/usr/local/bin

set -e
DOTFILES_VERSION=$(curl -sSL https://api.github.com/repos/fatindeed/dotfiles-go/tags |jq -r '.[0].name')
curl -sSLo "${INSTALL_DIR}/dotfiles" "https://github.com/fatindeed/dotfiles-go/releases/download/${DOTFILES_VERSION}/dotfiles-linux-amd64"
chmod a+x "${INSTALL_DIR}/dotfiles"
