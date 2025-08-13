#!/bin/bash

chmod 400 ./docker/mongo/keyfile
chown 999:999 ./docker/mongo/keyfile
chmod +x ./docker/mongo/setup.sh

docker-compose up -d
docker exec mogi-mongo sh /data/setup.sh
