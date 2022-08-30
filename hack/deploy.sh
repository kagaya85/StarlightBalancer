#!/usr/bin/env bash
set -eu

dir=./deployment

kubectl apply -f "$dir"/prometheus/ns.yaml

make load-image

kubectl apply -f "$dir"/config

kubectl apply -f "$dir"/prometheus -f "$dir"/balancer-deploy.yaml -f "$dir"/services-deploy.yaml

# helm repo add chaos-mesh https://charts.chaos-mesh.org

kubectl create ns chaos-testing

kind load docker-image ghcr.io/chaos-mesh/chaos-mesh:v2.3.0 --name dev

kind load docker-image ghcr.io/chaos-mesh/chaos-daemon:v2.3.0 --name dev

helm install chaos-mesh chaos-mesh/chaos-mesh -n=chaos-testing --set chaosDaemon.runtime=containerd --set chaosDaemon.socketPath=/run/containerd/containerd.sock --version v2.3.0
# curl -sSL https://mirrors.chaos-mesh.org/v2.3.0/install.sh | bash