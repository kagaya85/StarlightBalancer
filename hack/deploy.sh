#!/usr/bin/env bash
set -eu

dir=./deployment

kubectl apply -f "$dir"/prometheus/ns.yaml

make load-image

kubectl apply -f "$dir"/config

kubectl apply -f "$dir"/prometheus -f "$dir"/balancer-deploy.yaml -f "$dir"/services-deploy.yaml

helm install chaos-mesh chaos-mesh/chaos-mesh -n=chaos-testing --set chaosDaemon.runtime=containerd --set chaosDaemon.socketPath=/run/containerd/containerd.sock
# curl -sSL https://mirrors.chaos-mesh.org/v2.3.0/install.sh | bash