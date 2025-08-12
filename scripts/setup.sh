#!/bin/bash

chmod 400 ./docker/mongo/mongodb-keyfile
chown 999:999 ./docker/mongo/mongodb-keyfile
chmod +x ./docker/mongo/setup.sh
