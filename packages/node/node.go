package node

import (
	"fmt"
	"log"

	"github.com/nawesan12/fernet-token/packages/blockchain"
	"github.com/nawesan12/fernet-token/packages/p2p"
)

type Config struct {
	DataDir string
	P2PPort string
}

type Node struct {
	Blockchain *blockchain.Blockchain
	Mempool    *Mempool
	P2P        *p2p.P2PServer
	store      blockchain.Storage
	config     Config
}

func NewNode(cfg Config) (*Node, error) {
	store, err := blockchain.NewBoltStorage(cfg.DataDir + "/blockchain.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	bc, err := blockchain.NewBlockchain(store)
	if err != nil {
		store.Close()
		return nil, fmt.Errorf("failed to create blockchain: %w", err)
	}

	n := &Node{
		Blockchain: bc,
		Mempool:    NewMempool(),
		store:      store,
		config:     cfg,
	}

	n.P2P = p2p.NewP2PServer(cfg.P2PPort, n.handleP2PMessage)

	return n, nil
}

// NewNodeWithStorage creates a node with a custom storage (for testing).
func NewNodeWithStorage(store blockchain.Storage, p2pPort string) (*Node, error) {
	bc, err := blockchain.NewBlockchain(store)
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain: %w", err)
	}

	n := &Node{
		Blockchain: bc,
		Mempool:    NewMempool(),
		store:      store,
	}

	n.P2P = p2p.NewP2PServer(p2pPort, n.handleP2PMessage)

	return n, nil
}

// StartP2P starts the P2P server.
func (n *Node) StartP2P() {
	go n.P2P.Start()
}

// SubmitTransaction validates a transaction, adds it to the mempool, and broadcasts it.
func (n *Node) SubmitTransaction(tx *blockchain.Transaction) error {
	if err := n.Blockchain.ValidateTransaction(tx); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	n.Mempool.Add(tx)
	n.P2P.BroadcastTransaction(tx)
	log.Printf("Transaction %s submitted and broadcast", tx.ID)
	return nil
}

// Mine pulls transactions from the mempool, mines a block, and broadcasts it.
func (n *Node) Mine(miner string) (*blockchain.Block, error) {
	pending := n.Mempool.GetPending(blockchain.MaxTxPerBlock)
	block, err := n.Blockchain.MineBlock(miner, pending)
	if err != nil {
		return nil, err
	}

	// Remove confirmed transactions from mempool
	n.Mempool.RemoveConfirmed(block.Transactions)

	// Broadcast the new block
	n.P2P.BroadcastBlock(block)

	return block, nil
}

// handleP2PMessage routes incoming P2P messages.
func (n *Node) handleP2PMessage(msg p2p.Message) {
	switch msg.Type {
	case p2p.MsgTransaction:
		if msg.Transaction != nil {
			if err := n.Blockchain.ValidateTransaction(msg.Transaction); err != nil {
				log.Printf("Received invalid transaction: %v", err)
				return
			}
			n.Mempool.Add(msg.Transaction)
			log.Printf("Received transaction %s from peer", msg.Transaction.ID)
		}

	case p2p.MsgBlock:
		if msg.Block != nil {
			if err := n.Blockchain.AddBlock(msg.Block); err != nil {
				log.Printf("Received invalid block: %v", err)
				return
			}
			// Remove confirmed transactions from mempool
			n.Mempool.RemoveConfirmed(msg.Block.Transactions)
			log.Printf("Received and added block %d from peer", msg.Block.Index)
		}

	case p2p.MsgGetBlocks:
		chain := n.Blockchain.GetChain()
		n.P2P.SendChain(msg.SenderAddr, chain)
		log.Printf("Sent chain (%d blocks) to peer %s", len(chain), msg.SenderAddr)

	case p2p.MsgChain:
		if msg.Chain != nil && n.Blockchain.ShouldReplaceChain(msg.Chain) {
			if err := n.Blockchain.ReplaceChain(msg.Chain); err != nil {
				log.Printf("Failed to replace chain: %v", err)
				return
			}
			log.Printf("Replaced chain with longer chain (%d blocks)", len(msg.Chain))
		}

	case p2p.MsgPing:
		// Respond with pong - handled at P2P layer

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// Close shuts down the node.
func (n *Node) Close() error {
	n.P2P.Stop()
	return n.store.Close()
}
