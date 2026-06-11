#!/usr/bin/env bash

function install_jumpbox_prerequisites() {
  echo "Installing jumpbox prerequisites..."

  sudo apt-get update
  sudo apt-get install -y wget curl vim openssl

  echo "Done."
}

function download_binaries() {
  echo "Downloading binaries..."

  mkdir -p downloads
  wget -q --show-progress --https-only --timestamping \
    -P downloads \
    -i downloads-amd64.txt

  echo "Done."
}
