package p2p

import (
	"net"
	"testing"

	"github.com/nawesan12/fernet-token/packages/blockchain"
)

func TestMessageFraming(t *testing.T) {
	// Create a pipe to simulate a network connection
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Test transaction message
	tx := &blockchain.Transaction{
		ID:       "test-tx-id",
		Sender:   "sender-addr",
		Receiver: "receiver-addr",
		Amount:   100,
		Fee:      1,
	}

	sent := Message{
		Type:        MsgTransaction,
		Transaction: tx,
	}

	// Write in goroutine
	go func() {
		if err := WriteMessage(client, sent); err != nil {
			t.Errorf("WriteMessage failed: %v", err)
		}
	}()

	// Read
	received, err := ReadMessage(server)
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	if received.Type != MsgTransaction {
		t.Errorf("expected type %s, got %s", MsgTransaction, received.Type)
	}
	if received.Transaction == nil {
		t.Fatal("transaction should not be nil")
	}
	if received.Transaction.ID != tx.ID {
		t.Errorf("expected tx ID %s, got %s", tx.ID, received.Transaction.ID)
	}
	if received.Transaction.Amount != tx.Amount {
		t.Errorf("expected amount %d, got %d", tx.Amount, received.Transaction.Amount)
	}
}

func TestBlockMessage(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	block := &blockchain.Block{
		Index:    5,
		Hash:     "0000abc123",
		PrevHash: "0000def456",
		Miner:    "miner1",
	}

	sent := Message{
		Type:  MsgBlock,
		Block: block,
	}

	go func() {
		WriteMessage(client, sent)
	}()

	received, err := ReadMessage(server)
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	if received.Block == nil {
		t.Fatal("block should not be nil")
	}
	if received.Block.Index != 5 {
		t.Errorf("expected block index 5, got %d", received.Block.Index)
	}
	if received.Block.Hash != "0000abc123" {
		t.Errorf("expected hash '0000abc123', got '%s'", received.Block.Hash)
	}
}

func TestPingPong(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		WriteMessage(client, Message{Type: MsgPing})
	}()

	received, err := ReadMessage(server)
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	if received.Type != MsgPing {
		t.Errorf("expected PING, got %s", received.Type)
	}
}
