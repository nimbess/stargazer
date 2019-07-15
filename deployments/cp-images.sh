#!/bin/bash

# mntdir is the cifs mount from the kubernetes host to the build host
mntdir=shague

# hosts is an array of all the kubernetes hosts, e.g. master node-1 node-2
hosts=("192.168.122.73" "192.168.122.137" "192.168.122.165")

# virthost is the host containing the kubernetes hosts
virthost=172.31.1.50

echo "Executing docker save -o /tmp/stargazer.1.tar nimbess/stargazer:1"
docker save -o /tmp/stargazer.1.tar nimbess/stargazer:1

for ip in ${hosts[@]}; do
    cmd=(ssh -i ~/.ssh/vmhost/id_vm_rsa -o ProxyCommand="ssh -W %h:%p root@$virthost" centos@$ip "docker load -i /mnt/$mntdir/tmp/stargazer.1.tar")
    echo "Executing: ${cmd[@]}"
    "${cmd[@]}"
done
