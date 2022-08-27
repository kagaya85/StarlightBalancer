#!/usr/bin/env bash
set -eu

dir=./deployment

kubectl apply -f "$dir"/prometheus/ns.yaml

make load-image

kubectl apply -f "$dir"/config

kubectl apply -f "$dir"/prometheus -f "$dir"/balancer-deploy.yaml -f "$dir"/services-deploy.yaml

curl -sSL https://mirrors.chaos-mesh.org/v2.3.0/install.sh | bash