#!/bin/bash
# Docker for linux does not support hosts.docker.internal so we need to manually create the entry
# https://github.com/docker/for-linux/issues/264
# This fix is based on a workaround fix in the ticket, but it mysteriously uses a different hostname, and
# during testing we found that the default gateway did not work on mac.
CONTAINER_HOST_DNS=host.docker.internal

#https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds
#https://docs.docker.com/docker-for-windows/networking/#use-cases-and-workarounds

getent hosts $CONTAINER_HOST_DNS

if [ $? -ne 0 ]; then
  echo "Could not resolve $CONTAINER_HOST_DNS, manually setting to default gateway"

  echo -e "`/sbin/ip route|awk '/default/ { print $3 }'`\t$CONTAINER_HOST_DNS" | tee -a /etc/hosts > /dev/null
fi

#TODO: This sleep command is a naive implementation of enforcing order in buddy (we want postgres to start before AM)
#TODO: We rather need to implement waiting/reconnect logic to the personal data app instead of this sleep command
sleep 5
mongo_dsn=${MONGO_DSN?MONGO_DSN environment variable must be set}

set -eo pipefail

run_migrations() {
  echo "Starting Migration..."

  /migrate -verbose \
    -database "$mongo_dsn" \
    -path /db/migrations \
    up

  echo "Migrations complete..."
}

if [ -z "$RUN_MIGRATIONS" ]; then
  echo "RUN_MIGRATIONS environment variable must be set"
  exit 1
fi

case "$RUN_MIGRATIONS" in
  true)
    run_migrations
    ;;
  only)
    run_migrations
    echo "RUN_MIGRATIONS set to $RUN_MIGRATIONS exiting..."
    exit 0
    ;;
  false)
    ;;
  *)
    echo "Unknown value for RUN_MIGRATIONS: $RUN_MIGRATIONS"
    exit 2
    ;;
esac



exec /personal-data

