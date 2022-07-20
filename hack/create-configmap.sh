#!/usr/bin/env bash
set -eux


echo "create balancer configmap"
kubectl create configmap balancer-config --from-file=balancer/configs/config.yaml -o yaml --dry-run=client > deployment/config/balancer-configmap.yaml


echo "create service configmap"
for svc in $(ls services); do
    dir=services/"${svc}"
	if [[ -d $dir ]]; then
		kubectl create configmap "${svc}"-config --from-file="${dir}"/configs/config.yaml -o yaml --dry-run=client > deployment/config/"${svc}"-configmap.yaml
	fi
done