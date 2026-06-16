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

function extract_downloads() {
  echo "Extracting downloads..."

  mkdir -p downloads/{client,controller,worker}
  tar xvf downloads/crictl-v1.36.0-linux-amd64.tar.gz \
    -C downloads/worker/
  tar xvf downloads/containerd-2.3.1-linux-amd64.tar.gz  \
    --strip-components 1 \
    -C downloads/worker/
  tar xvf downloads/etcd-v3.6.12-linux-amd64.tar.gz \
    --strip-components 1 \
    -C downloads/ \
    etcd-v3.6.12-linux-amd64/etcdctl \
    etcd-v3.6.12-linux-amd64/etcd

  mv downloads/{etcdctl,kubectl} downloads/client/
  mv downloads/{etcd,kube-apiserver,kube-controller-manager,kube-scheduler} \
    downloads/controller/
  mv downloads/{kubelet,kube-proxy} downloads/worker/
  mv downloads/runc.amd64 downloads/worker/runc

  chmod +x downloads/{client,controller,worker}/*

  echo "Deleting archives..."
  rm -rf downloads/*.gz

  echo "Done."
}

function generate_certificates() {
  echo "Generating CA key and certificate..."
  
  openssl genrsa -out ca.key 4096
  openssl req -x509 -new -sha512 -noenc \
    -key ca.key -days 3653 \
    -subj "/CN=KUBERNETES-CA" \
    -addext "basicConstraints = critical,CA:TRUE" \
    -out ca.crt
}
