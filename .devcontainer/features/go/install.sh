#!/bin/bash
set -e

GO_VERSION=1.25.0

echo "Installing Go..."

curl -LO https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile.d/go.sh
go env -w GOPATH=/home/vscode/go
echo "export PATH=$PATH:$(go env GOPATH)/bin" >> /etc/profile.d/go.sh

GO_TOOLS=(
    "golang.org/x/tools/gopls@latest"
    "golang.org/x/tools/cmd/goimports@latest"
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    "google.golang.org/protobuf/cmd/protoc-gen-go@latest"
)

for tool in "${GO_TOOLS[@]}"; do
    echo "- $tool"
    go install -v "$tool"
    if [ $? -eq 0 ]; then
        echo "-> $tool installed."
    else
        echo "-> $tool installation failed. Please check the error."
        exit 1
    fi
done
