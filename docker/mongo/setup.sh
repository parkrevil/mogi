#!/bin/bash

echo "Waiting for MongoDB to be ready..."
sleep 10

echo "Setting up Single Node Replica Set..."
mongosh --host mongo:27017 -u mogi -p 1234 --authenticationDatabase admin --eval "
rs.initiate({
  _id: 'rs0',
  members: [
    {_id: 0, host: 'localhost:27017'}
  ]
})
"

echo "Waiting for replica set to be ready..."
sleep 10

echo "Creating database and collections..."
mongosh --host mongo:27017 -u mogi -p 1234 --authenticationDatabase admin --eval "
use mogi
db.createCollection('mogi')
"

echo "MongoDB single node replica set setup completed!"
