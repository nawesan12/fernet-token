package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
	Miner        string
}

type Blockchain struct {
	Chain       []Block
	PendingTxns []Transaction
	Balances    map[string]float64 // Track balances per address
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Balances: make(map[string]float64),
	}
	bc.CreateGenesisBlock()
	return bc
}

func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Hash:         "genesis_hash",
		Nonce:        0,
		Miner:        "",
	}
	bc.Chain = append(bc.Chain, genesisBlock)
	bc.Balances["network"] = 1000000
}

func (bc *Blockchain) AddTransaction(sender, receiver string, amount float64) bool {
	if sender == receiver || amount <= 0 {
		fmt.Println("❌ Invalid transaction: sender and receiver must be different, and amount must be positive.")
		return false
	}

	if bc.Balances[sender] < amount {
		fmt.Println("❌ Transaction failed: insufficient balance")
		return false
	}

	tx := NewTransaction(sender, receiver, amount)
	bc.PendingTxns = append(bc.PendingTxns, *tx)
	fmt.Println("📩 Transaction added:", tx)
	return true
}

func (bc *Blockchain) GetPendingTransactions() []Transaction {
	return bc.PendingTxns
}

func (bc *Blockchain) MineBlock(miner string, reward float64) {
	if len(bc.PendingTxns) == 0 {
		fmt.Println("⛏ No transactions to mine.")
		return
	}

	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:        len(bc.Chain),
		Timestamp:    time.Now().Unix(),
		Transactions: bc.PendingTxns,
		PrevHash:     prevBlock.Hash,
		Nonce:        0,
		Miner:        miner,
	}

	newBlock.Hash = bc.ProofOfWork(&newBlock)
	bc.Chain = append(bc.Chain, newBlock)
	bc.PendingTxns = nil

	for _, tx := range newBlock.Transactions {
		bc.Balances[tx.Sender] -= tx.Amount
		bc.Balances[tx.Receiver] += tx.Amount
	}

	bc.Balances[miner] += reward

	fmt.Println("✅ Block mined by", miner, "with hash:", newBlock.Hash)
}

func (bc *Blockchain) ProofOfWork(block *Block) string {
	var hash string
	for {
		block.Nonce++
		hash = bc.CalculateHash(block)
		if hash[:4] == "0000" {
			return hash
		}
	}
}

func (bc *Blockchain) CalculateHash(block *Block) string {
	blockData, _ := json.Marshal(block)
	hash := sha256.Sum256(blockData)
	return fmt.Sprintf("%x", hash)
}

func (bc *Blockchain) GetBalance(address string) float64 {
	return bc.Balances[address]
}

// NewTransaction creates a new signed transaction
type Consensus interface {
	ValidateTransaction(tx *Transaction) bool
}
