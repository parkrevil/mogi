#!/bin/bash

set -e

echo "🚀 Setting up Mogi Development Environment..."
echo "📦 Setting up system..."
echo "[SYSTEM] Updating packages..."
sudo apt-get update -qq -y

echo "[SYSTEM] Installing packages..."
sudo apt-get install -y -qq \
    protobuf-compiler

echo "[SYSTEM] Upgrading packages..."
sudo apt-get upgrade -qq -y

sudo mkdir -p /home/vscode/go
sudo chown -R vscode:vscode /home/vscode/go

sudo mkdir -p /home/vscode/.cache/go-build
sudo chown -R vscode:vscode /home/vscode/.cache

echo "🐳 Setting up Docker"

NAMESPACE="mogi-dev"
SERVICES=(mongo redis-stack)

echo "[DOCKER] Checking availability..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found."
    exit 1
fi

echo "[DOCKER] Checking daemon..."
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker daemon is not running."
    exit 1
fi

echo "[DOCKER] Building images..."
for service in "${SERVICES[@]}"; do
    docker build -f docker/$service/Dockerfile -t $NAMESPACE-$service docker/$service
done

echo "[DOCKER] Starting containers..."

mkdir -p ./.container-volumes/mongo/data
mkdir -p ./.container-volumes/redis-stack/data
docker compose -f docker/docker-compose.yml up -d

echo "[MONGODB] Waiting for startup..."
until docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  echo "[MONGODB] ⏳ MongoDB is not ready yet, waiting..."
  sleep 2
done

echo "[MONGODB] Setting up replica set..."
docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "
try {
  rs.initiate({
    _id: 'rs0',
    members: [
      {_id: 0, host: 'localhost:27117'}
    ]
  });
  print('Replica set initialized successfully');
} catch (error) {
  if (error.message.includes('already initialized') || error.message.includes('already a member')) {
    print('Replica set already initialized');
  } else {
    print('Error setting up replica set: ' + error.message);
  }
}
"

echo "[MONGODB] Waiting for replica set to be ready..."
until docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "rs.status()" > /dev/null 2>&1; do
  echo "[MONGODB] ⏳ Replica set is not ready yet, waiting..."
  sleep 2
done

echo "[MONGODB] Creating database and collections..."
docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "
use mogi
db.createCollection('mogi')
"

echo "[REDIS-STACK] Waiting for startup..."
until docker exec mogi-dev-redis-stack redis-cli -u redis://mogi:1234@localhost:6379 ping > /dev/null 2>&1; do
  echo "[REDIS-STACK] ⏳ Redis Stack is not ready yet, waiting..."
  sleep 2
done

echo "🔧 Install project dependencies..."
bun install
go work sync

echo ""
echo "🎉 Mogi Development Environment setup completed!"
echo "🚀 Ready to start development!"
