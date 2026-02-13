package p2p

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/nawesan12/fernet-token/packages/blockchain"
)

const (
	MsgTransaction = "TRANSACTION"
	MsgBlock       = "BLOCK"
	MsgGetBlocks   = "GET_BLOCKS"
	MsgChain       = "CHAIN"
	MsgPing        = "PING"
	MsgPong        = "PONG"

	MaxMessageSize = 10 * 1024 * 1024 // 10MB
)

// Message is the wire format for P2P communication.
type Message struct {
	Type        string                    `json:"type"`
	Transaction *blockchain.Transaction   `json:"transaction,omitempty"`
	Block       *blockchain.Block         `json:"block,omitempty"`
	Chain       []blockchain.Block        `json:"chain,omitempty"`
	SenderAddr  string                    `json:"senderAddr,omitempty"`
}

// WriteMessage writes a length-prefixed JSON message to a connection.
func WriteMessage(conn net.Conn, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write 4-byte big-endian length prefix
	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	// Write payload
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
}

// ReadMessage reads a length-prefixed JSON message from a connection.
func ReadMessage(conn net.Conn) (Message, error) {
	var msg Message

	// Read 4-byte big-endian length prefix
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return msg, fmt.Errorf("failed to read length: %w", err)
	}

	if length > MaxMessageSize {
		return msg, fmt.Errorf("message too large: %d bytes", length)
	}

	// Read exact payload
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return msg, fmt.Errorf("failed to read payload: %w", err)
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return msg, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return msg, nil
}
