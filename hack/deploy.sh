#!/usr/bin/env bash
set -eu

dir=./deployment

kubectl apply -f "$dir"/prometheus/ns.yaml

make load-image

kubectl apply -f "$dir"/config

kubectl apply -f "$dir"/prometheus -f "$dir"/balancer-deploy.yaml -f "$dir"/services-deploy.yaml