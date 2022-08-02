#!/usr/bin/env bash
set -eu

dir=./deployment

k apply -f "$dir"/prometheus/ns.yaml

make load-image

k apply -f "$dir"/config

k apply -f "$dir"/prometheus -f "$dir"/balncer-deploy.yaml -f "$dir"/services-deploy.yaml