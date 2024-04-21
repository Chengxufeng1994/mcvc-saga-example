#!/bin/bash

REDIES_NODES=$(seq 1 6)
REDIS_CONFG_DIR=$(PWD)/infra/redis
REDIS_CONFG_FILE=redis.conf
REDIS_CLUSTER_IP=$(ifconfig | grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" \
    | grep -v 127.0.0.1 | awk '{ print $2 }' | cut -f2 -d: | head -n1)

for node in ${REDIES_NODES};
do
  port=700${node}
  mkdir -p ${REDIS_CONFG_DIR}/node-$node/conf
  touch ${REDIS_CONFG_DIR}/node-$node/conf/${REDIS_CONFG_FILE}
  cat << EOF > ${REDIS_CONFG_DIR}/node-$node/conf/${REDIS_CONFG_FILE}
bind 0.0.0.0
port ${port}
appendonly yes
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
cluster-announce-ip ${REDIS_CLUSTER_IP}
cluster-announce-port ${port}
cluster-announce-bus-port 1${port}
masterauth password
requirepass password
EOF
done