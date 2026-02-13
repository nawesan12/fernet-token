#!/bin/bash
set -e

echo "=== Fernet Token Network Test ==="
echo "Starting two nodes and verifying sync..."

TMPDIR=$(mktemp -d)
trap "kill 0; rm -rf $TMPDIR" EXIT

# Start Node A
echo "Starting Node A (HTTP:8080, P2P:6000)..."
go run ./apps/api \
    --http-port 8080 \
    --p2p-port 6000 \
    --data-dir "$TMPDIR/nodeA" &
NODE_A_PID=$!

sleep 2

# Start Node B, connecting to Node A
echo "Starting Node B (HTTP:8081, P2P:6001, peer:localhost:6000)..."
go run ./apps/api \
    --http-port 8081 \
    --p2p-port 6001 \
    --data-dir "$TMPDIR/nodeB" \
    --peers "localhost:6000" &
NODE_B_PID=$!

sleep 2

echo ""
echo "--- Creating wallet on Node A ---"
WALLET=$(curl -s -X POST http://localhost:8080/api/wallet/create)
echo "Wallet: $WALLET"
ADDRESS=$(echo "$WALLET" | grep -o '"address":"[^"]*"' | cut -d'"' -f4)
echo "Address: $ADDRESS"

echo ""
echo "--- Requesting faucet on Node A ---"
curl -s -X POST http://localhost:8080/api/faucet \
    -H "Content-Type: application/json" \
    -d "{\"address\": \"$ADDRESS\"}"
echo ""

echo ""
echo "--- Mining block on Node A ---"
curl -s -X POST http://localhost:8080/api/mine \
    -H "Content-Type: application/json" \
    -d "{\"miner\": \"$ADDRESS\"}"
echo ""

sleep 2

echo ""
echo "--- Checking Node A height ---"
curl -s http://localhost:8080/api/blockchain/height
echo ""

echo ""
echo "--- Checking Node B height ---"
curl -s http://localhost:8081/api/blockchain/height
echo ""

echo ""
echo "--- Checking balance on Node A ---"
curl -s "http://localhost:8080/api/balance/$ADDRESS"
echo ""

echo ""
echo "=== Test complete ==="
