#!/usr/bin/env bash
set -eux

echo
echo "Start build balancer image, Repo: $1, Tag: $2"
echo

docker build -f balancer/Dockerfile -t "$1"/starlight-balancer:"$2" .

echo
echo "Start build service images, Repo: $1, Tag: $2"
echo
for svc in $(ls services); do
    dir=services/"${svc}"
    if [[ -d $dir ]]; then
        if [[ -n $(ls "$dir" | grep -i Dockerfile) ]]; then
            echo "build ${dir}"
            docker build -f "${dir}"/Dockerfile -t "$1"/starlight-"${svc}":"$2" .
        fi
    fi
done