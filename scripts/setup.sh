#!/bin/bash
set -e

echo "Setting up Fernet Token node..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.22+ from https://go.dev/dl/"
    exit 1
fi

echo "Go version: $(go version)"

# Download Go dependencies
echo "Downloading dependencies..."
go mod tidy

# Set defaults
HTTP_PORT=${HTTP_PORT:-8080}
P2P_PORT=${P2P_PORT:-6000}
DATA_DIR=${DATA_DIR:-"$HOME/.fernet-token"}

mkdir -p "$DATA_DIR"
echo "Data directory: $DATA_DIR"

# Start the node
echo "Starting Fernet Token node..."
echo "  HTTP API: http://localhost:$HTTP_PORT"
echo "  P2P Port: $P2P_PORT"

go run ./apps/api \
    --http-port "$HTTP_PORT" \
    --p2p-port "$P2P_PORT" \
    --data-dir "$DATA_DIR"
