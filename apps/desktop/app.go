package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nawesan12/fernet-token/packages/blockchain"
	"github.com/nawesan12/fernet-token/packages/node"
	"github.com/nawesan12/fernet-token/packages/wallet"
)

type App struct {
	ctx    context.Context
	node   *node.Node
	wallet *wallet.Wallet
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".fernet-token-desktop")
	os.MkdirAll(dataDir, 0755)

	cfg := node.Config{
		DataDir: dataDir,
		P2PPort: "6100",
	}

	n, err := node.NewNode(cfg)
	if err != nil {
		log.Printf("Failed to create node: %v", err)
		return
	}
	a.node = n
	n.StartP2P()

	// Load or create wallet
	keyPath := filepath.Join(dataDir, "wallet.pem")
	w, err := wallet.LoadFromFile(keyPath)
	if err != nil {
		w, err = wallet.NewWallet()
		if err != nil {
			log.Printf("Failed to create wallet: %v", err)
			return
		}
		w.SaveToFile(keyPath)
	}
	a.wallet = w
	log.Printf("Wallet loaded: %s", w.Address)
}

func (a *App) shutdown(ctx context.Context) {
	if a.node != nil {
		a.node.Close()
	}
}

// CreateWallet generates a new wallet (replaces current).
func (a *App) CreateWallet() map[string]string {
	home, _ := os.UserHomeDir()
	keyPath := filepath.Join(home, ".fernet-token-desktop", "wallet.pem")

	w, err := wallet.NewWallet()
	if err != nil {
		return map[string]string{"error": err.Error()}
	}
	w.SaveToFile(keyPath)
	a.wallet = w

	return map[string]string{
		"address":   w.Address,
		"publicKey": w.PublicKey,
	}
}

// GetWalletInfo returns current wallet details.
func (a *App) GetWalletInfo() map[string]string {
	if a.wallet == nil {
		return map[string]string{"error": "no wallet"}
	}
	return map[string]string{
		"address":   a.wallet.Address,
		"publicKey": a.wallet.PublicKey,
	}
}

// GetBalance returns the balance of the current wallet.
func (a *App) GetBalance() map[string]interface{} {
	if a.wallet == nil || a.node == nil {
		return map[string]interface{}{"balance": 0, "formatted": "0 FERNET"}
	}
	bal := a.node.Blockchain.GetBalance(a.wallet.Address)
	return map[string]interface{}{
		"balance":   bal,
		"formatted": formatFernet(bal),
	}
}

// SendTransaction signs and submits a transaction.
func (a *App) SendTransaction(receiver string, amount, fee uint64) map[string]string {
	if a.wallet == nil || a.node == nil {
		return map[string]string{"error": "node not initialized"}
	}

	nonce := a.node.Blockchain.GetNonce(a.wallet.Address)
	tx := blockchain.NewTransaction(a.wallet.Address, receiver, amount, fee, nonce, a.wallet.PublicKey)

	sig, err := a.wallet.Sign(tx.SignableData())
	if err != nil {
		return map[string]string{"error": fmt.Sprintf("signing failed: %v", err)}
	}
	tx.Signature = sig

	if err := a.node.SubmitTransaction(tx); err != nil {
		return map[string]string{"error": fmt.Sprintf("submit failed: %v", err)}
	}

	return map[string]string{
		"message": "transaction submitted",
		"txId":    tx.ID,
	}
}

// Mine triggers mining.
func (a *App) Mine() map[string]interface{} {
	if a.wallet == nil || a.node == nil {
		return map[string]interface{}{"error": "node not initialized"}
	}

	block, err := a.node.Mine(a.wallet.Address)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	return map[string]interface{}{
		"message": fmt.Sprintf("Mined block #%d", block.Index),
		"index":   block.Index,
		"hash":    block.Hash,
	}
}

// GetChainHeight returns the blockchain height.
func (a *App) GetChainHeight() uint64 {
	if a.node == nil {
		return 0
	}
	return a.node.Blockchain.Height()
}

// GetBlock returns a block by index.
func (a *App) GetBlock(index uint64) map[string]interface{} {
	if a.node == nil {
		return map[string]interface{}{"error": "node not initialized"}
	}
	block, err := a.node.Blockchain.GetBlock(index)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	return map[string]interface{}{
		"index":        block.Index,
		"hash":         block.Hash,
		"prevHash":     block.PrevHash,
		"timestamp":    block.Timestamp,
		"miner":        block.Miner,
		"nonce":        block.Nonce,
		"transactions": block.Transactions,
	}
}

// GetRecentBlocks returns the last N blocks.
func (a *App) GetRecentBlocks(count int) []map[string]interface{} {
	if a.node == nil {
		return nil
	}
	chain := a.node.Blockchain.GetChain()
	start := len(chain) - count
	if start < 0 {
		start = 0
	}

	var result []map[string]interface{}
	for i := len(chain) - 1; i >= start; i-- {
		b := chain[i]
		result = append(result, map[string]interface{}{
			"index":     b.Index,
			"hash":      b.Hash,
			"timestamp": b.Timestamp,
			"miner":     b.Miner,
			"txCount":   len(b.Transactions),
		})
	}
	return result
}

// ConnectToPeer connects to a remote peer.
func (a *App) ConnectToPeer(address string) map[string]string {
	if a.node == nil {
		return map[string]string{"error": "node not initialized"}
	}
	if err := a.node.P2P.ConnectToPeer(address); err != nil {
		return map[string]string{"error": err.Error()}
	}
	return map[string]string{"message": "connected to " + address}
}

// GetPeerCount returns the number of connected peers.
func (a *App) GetPeerCount() int {
	if a.node == nil {
		return 0
	}
	return a.node.P2P.PeerCount()
}

// GetPeers returns addresses of connected peers.
func (a *App) GetPeers() []string {
	if a.node == nil {
		return nil
	}
	return a.node.P2P.PeerAddresses()
}

// GetPendingTxCount returns mempool size.
func (a *App) GetPendingTxCount() int {
	if a.node == nil {
		return 0
	}
	return a.node.Mempool.Count()
}

func formatFernet(fernetoshi uint64) string {
	whole := fernetoshi / blockchain.OneFernet
	frac := fernetoshi % blockchain.OneFernet
	if frac == 0 {
		return fmt.Sprintf("%d FERNET", whole)
	}
	return fmt.Sprintf("%d.%08d FERNET", whole, frac)
}
