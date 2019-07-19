#!/bin/bash

if [ "$1" == "-h" ]; then
    echo "$0 [apply|delete]"
    exit
fi

operation=${1:-apply}
specs=("etcd" "serviceaccount" "rbac" "rbac-cluster" "pod")

# mntdir is the path to where the stargazer repo is located
mntdir="/mnt/shague/go/src/github.com/nimbess"

for spec in ${specs[@]}; do
    cmd=(kubectl "$operation" -f "$mntdir/stargazer/deployments/stargazer-$spec.yaml")
    echo "Executing: ${cmd[@]}"
    "${cmd[@]}"
done
