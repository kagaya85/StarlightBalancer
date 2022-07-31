#!/usr/bin/env bash
set -eu

echo
echo "Load images to kind, Repo: $1"
echo
images=$(docker images | grep "$1"/starlight- | awk 'BEGIN{OFS=":"}{print $1,$2}')

if [[ -n "$images" ]]; then
    echo "$images" | xargs -I {} kind load docker-image {} --name dev
fi