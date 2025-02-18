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
}

// NewBlockchain crea una nueva instancia de Blockchain
func NewBlockchain() *Blockchain {
	bc := &Blockchain{}
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
}

func (bc *Blockchain) AddTransaction(sender, receiver string, amount float64) {
	tx := Transaction{Sender: sender, Receiver: receiver, Amount: amount}
	bc.PendingTxns = append(bc.PendingTxns, tx)
	fmt.Println("Transaction added:", tx)
}

func (bc *Blockchain) GetPendingTransactions() []Transaction {
	return bc.PendingTxns
}

func (bc *Blockchain) MineBlock(miner string, reward float64) {
	if len(bc.PendingTxns) == 0 {
		fmt.Println("No transactions to mine")
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

	bc.AddTransaction("network", miner, reward)

	fmt.Println("Block mined by: ", miner, "with hash:", newBlock.Hash)
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
