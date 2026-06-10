#!/usr/bin/env bash

function install_jumpbox_prerequisites() {
  echo "Installing jumpbox prerequisites..."

  sudo apt-get update
  sudo apt-get install -y wget curl vim openssl git

  echo "Done."
}
