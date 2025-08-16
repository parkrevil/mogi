#!/bin/bash

set -e

echo "🚀 Setting up Mogi Development Environment..."

# =============================================================================
# 1. SYSTEM SETUP
# =============================================================================
echo "📦 Updating system packages..."
sudo apt-get update -qq

echo "🛠️ Installing essential utilities..."
sudo apt-get install -y -qq \
    curl \
    wget \
    git \
    vim \
    htop \
    tree \
    jq \
    unzip \
    protobuf-compiler

echo "📦 Upgrading system packages..."
sudo apt-get upgrade -qq

# =============================================================================
# 2. DOCKER SETUP
# =============================================================================
echo "🐳 Setting up Docker environment..."

# Check if Docker is available
echo "🐳 Checking Docker availability..."
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose not found. Please ensure Docker Desktop is properly installed."
    exit 1
fi

# Check if Docker daemon is running
echo "🔍 Checking Docker daemon..."
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker daemon is not running. Please:"
    echo "   1. Start Docker Desktop on your host machine"
    echo "   2. Ensure Docker Desktop is running and accessible"
    echo "   3. Restart your Dev Container"
    exit 1
fi

# Verify installations
echo "🔍 Verifying Docker installations..."
docker --version
docker-compose --version
echo "✅ Docker and Docker Compose installed"

# Set Docker BuildKit to 0 to avoid bake definition issues
echo "🔧 Setting Docker BuildKit to 0..."
echo 'export DOCKER_BUILDKIT=0' >> ~/.bashrc
export DOCKER_BUILDKIT=0
source ~/.bashrc

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p ./.container-volumes/mongo/data
mkdir -p ./.container-volumes/redis-stack/data

# =============================================================================
# 3. CONTAINERS SETUP
# =============================================================================
# Build and start containers
echo "🔨 Building containers..."
docker compose -f ./docker-compose.yml build --no-cache

echo "🚀 Starting containers..."
docker compose -f ./docker-compose.yml up -d

# Wait for MongoDB to be ready
echo "⏳ Waiting for MongoDB to be ready..."
until docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  echo "  ⏳ MongoDB is not ready yet, waiting..."
  sleep 2
done
echo "✅ MongoDB is ready!"

# Setup MongoDB Replica Set
echo "🔄 Setting up MongoDB Replica Set..."
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

# Wait for replica set to be ready
echo "⏳ Waiting for replica set to be ready..."
until docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "rs.status()" > /dev/null 2>&1; do
  echo "  ⏳ Replica set is not ready yet, waiting..."
  sleep 2
done
echo "✅ Replica set is ready!"

# Create database and collections
echo "🗄️ Creating database and collections..."
docker exec mogi-dev-mongo mongosh --port 27117 -u mogi -p 1234 --authenticationDatabase admin --eval "
use mogi
db.createCollection('mogi')
"
echo "✅ MongoDB setup completed!"

# Wait for Redis Stack to be ready
echo "⏳ Waiting for Redis Stack to be ready..."
until docker exec mogi-dev-redis-stack redis-cli -u redis://mogi:1234@localhost:6379 ping > /dev/null 2>&1; do
  echo "  ⏳ Redis Stack is not ready yet, waiting..."
  sleep 2
done
echo "✅ Redis Stack is ready!"

# =============================================================================
# 4. GIT SETUP
# =============================================================================
echo "🔧 Setting up Git configuration..."

# Verify git configuration
if [ -f ~/.gitconfig ]; then
    echo "✅ Git configuration found"
    echo "👤 Git user: $(git config user.name)"
    echo "📧 Git email: $(git config user.email)"
else
    echo "⚠️  Git configuration not found. Please ensure ~/.gitconfig is mounted from host"
fi

# =============================================================================
# 5. DEVELOPMENT TOOLS SETUP
# =============================================================================
echo "🔧 Setting up development tools..."

# Install Bun
echo "🍞 Installing Bun v1.2.20..."
curl -fsSL https://bun.sh/install | bash -s "bun-v1.2.20"
export PATH="$HOME/.bun/bin:$PATH"
echo 'export PATH="$HOME/.bun/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Install Go tools
echo "🐹 Installing Go tools..."
# Golang 1.25.0 설치 (공식 문서 참고)
echo "🐹 Installing Go v1.25.0..."
GO_VERSION=1.25.0
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
  ARCH=amd64
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH=arm64
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi
wget -q https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz -O /tmp/go${GO_VERSION}.linux-${ARCH}.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-${ARCH}.tar.gz
export PATH="/usr/local/go/bin:$PATH"
echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.bashrc
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
go version
go install golang.org/x/tools/cmd/goimports@v0.36.0
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.7

# =============================================================================
# 6. PROJECT DEPENDENCIES
# =============================================================================
echo "📚 Installing project dependencies..."

# Clean existing node_modules
echo "🧹 Cleaning existing node_modules..."
rm -rf node_modules
rm -rf apps/*/node_modules
rm -rf packages/*/node_modules

# Install Bun dependencies
echo "📦 Installing Bun dependencies..."
bun install

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go work sync

# =============================================================================
# 7. COMPLETION
# =============================================================================
echo ""
echo "🎉 Mogi Development Environment setup completed!"
echo "🚀 Ready to start development!"
