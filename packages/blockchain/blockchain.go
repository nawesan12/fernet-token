package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Blockchain struct {
	Chain    []Block
	Balances map[string]uint64
	Nonces   map[string]uint64
	store    Storage
	mu       sync.RWMutex
}

// NewBlockchain creates or loads a blockchain from storage.
func NewBlockchain(store Storage) (*Blockchain, error) {
	bc := &Blockchain{
		Balances: make(map[string]uint64),
		Nonces:   make(map[string]uint64),
		store:    store,
	}

	chain, err := store.LoadChain()
	if err != nil {
		return nil, fmt.Errorf("failed to load chain: %w", err)
	}

	if len(chain) > 0 {
		bc.Chain = chain
		bc.rebuildState()
		log.Printf("Loaded blockchain with %d blocks from storage", len(bc.Chain))
	} else {
		bc.createGenesisBlock()
		log.Println("Created new blockchain with genesis block")
	}

	return bc, nil
}

func (bc *Blockchain) createGenesisBlock() {
	genesis := Block{
		Index:        0,
		Timestamp:    GenesisTimestamp,
		Transactions: []Transaction{},
		PrevHash:     "0",
		Nonce:        0,
		Miner:        "",
	}
	genesis.Hash = CalculateBlockHash(&genesis)
	bc.Chain = append(bc.Chain, genesis)
	bc.store.SaveBlock(genesis)
}

// rebuildState replays the chain to reconstruct balances and nonces.
func (bc *Blockchain) rebuildState() {
	bc.Balances = make(map[string]uint64)
	bc.Nonces = make(map[string]uint64)

	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.Sender == CoinbaseSender {
				bc.Balances[tx.Receiver] += tx.Amount
			} else {
				bc.Balances[tx.Sender] -= (tx.Amount + tx.Fee)
				bc.Balances[tx.Receiver] += tx.Amount
				bc.Balances[block.Miner] += tx.Fee
				if tx.Nonce >= bc.Nonces[tx.Sender] {
					bc.Nonces[tx.Sender] = tx.Nonce + 1
				}
			}
		}
	}
}

// CalculateBlockHash computes the hash of a block using BlockHashData (excludes Hash field).
func CalculateBlockHash(block *Block) string {
	data := BlockHashData{
		Index:        block.Index,
		Timestamp:    block.Timestamp,
		Transactions: block.Transactions,
		PrevHash:     block.PrevHash,
		Nonce:        block.Nonce,
		Miner:        block.Miner,
	}
	jsonBytes, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("%x", hash)
}

// MineBlock creates a new block with the given transactions.
func (bc *Blockchain) MineBlock(miner string, pendingTxns []Transaction) (*Block, error) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Validate all transactions
	var validTxns []Transaction
	for _, tx := range pendingTxns {
		if err := bc.validateTransactionLocked(&tx); err != nil {
			log.Printf("Skipping invalid tx %s: %v", tx.ID, err)
			continue
		}
		validTxns = append(validTxns, tx)
		if len(validTxns) >= MaxTxPerBlock {
			break
		}
	}

	// Add coinbase transaction
	coinbase := NewCoinbaseTx(miner, MiningReward)
	allTxns := append([]Transaction{*coinbase}, validTxns...)

	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: allTxns,
		PrevHash:     prevBlock.Hash,
		Nonce:        0,
		Miner:        miner,
	}

	// Proof of Work
	for {
		newBlock.Hash = CalculateBlockHash(&newBlock)
		if strings.HasPrefix(newBlock.Hash, TargetPrefix) {
			break
		}
		newBlock.Nonce++
	}

	// Apply state changes
	bc.Balances[miner] += MiningReward
	for _, tx := range validTxns {
		bc.Balances[tx.Sender] -= (tx.Amount + tx.Fee)
		bc.Balances[tx.Receiver] += tx.Amount
		bc.Balances[miner] += tx.Fee
		if tx.Nonce >= bc.Nonces[tx.Sender] {
			bc.Nonces[tx.Sender] = tx.Nonce + 1
		}
	}

	bc.Chain = append(bc.Chain, newBlock)

	// Persist
	bc.store.SaveBlock(newBlock)
	bc.store.SaveBalances(bc.Balances)
	bc.store.SaveNonces(bc.Nonces)

	log.Printf("Block %d mined by %s with hash %s (%d txns)", newBlock.Index, miner, newBlock.Hash, len(validTxns))
	return &newBlock, nil
}

// ValidateTransaction checks if a transaction is valid against current state.
func (bc *Blockchain) ValidateTransaction(tx *Transaction) error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.validateTransactionLocked(tx)
}

func (bc *Blockchain) validateTransactionLocked(tx *Transaction) error {
	if err := tx.IsValid(); err != nil {
		return err
	}

	if tx.Sender == CoinbaseSender {
		return nil
	}

	// Check balance
	balance := bc.Balances[tx.Sender]
	if balance < tx.Amount+tx.Fee {
		return fmt.Errorf("insufficient balance: has %d, needs %d", balance, tx.Amount+tx.Fee)
	}

	// Check nonce
	expectedNonce := bc.Nonces[tx.Sender]
	if tx.Nonce != expectedNonce {
		return fmt.Errorf("invalid nonce: expected %d, got %d", expectedNonce, tx.Nonce)
	}

	return nil
}

// ValidateChain verifies the entire chain integrity.
func (bc *Blockchain) ValidateChain() error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.Chain) == 0 {
		return errors.New("chain is empty")
	}

	// Check genesis
	if bc.Chain[0].PrevHash != "0" {
		return errors.New("invalid genesis block")
	}

	for i := 1; i < len(bc.Chain); i++ {
		block := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		// Check index sequence
		if block.Index != prevBlock.Index+1 {
			return fmt.Errorf("block %d: invalid index", i)
		}

		// Check prev hash link
		if block.PrevHash != prevBlock.Hash {
			return fmt.Errorf("block %d: prev hash mismatch", i)
		}

		// Check hash correctness
		expectedHash := CalculateBlockHash(&block)
		if block.Hash != expectedHash {
			return fmt.Errorf("block %d: hash mismatch", i)
		}

		// Check PoW
		if !strings.HasPrefix(block.Hash, TargetPrefix) {
			return fmt.Errorf("block %d: insufficient proof of work", i)
		}

		// Check coinbase
		if len(block.Transactions) == 0 || block.Transactions[0].Sender != CoinbaseSender {
			return fmt.Errorf("block %d: missing coinbase transaction", i)
		}
		if block.Transactions[0].Amount != MiningReward {
			return fmt.Errorf("block %d: invalid coinbase amount", i)
		}

		// Verify all non-coinbase transaction signatures
		for j := 1; j < len(block.Transactions); j++ {
			if err := block.Transactions[j].VerifySignature(); err != nil {
				return fmt.Errorf("block %d, tx %d: %w", i, j, err)
			}
		}
	}

	return nil
}

// GetBalance returns the balance for an address.
func (bc *Blockchain) GetBalance(address string) uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.Balances[address]
}

// GetNonce returns the next expected nonce for an address.
func (bc *Blockchain) GetNonce(address string) uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.Nonces[address]
}

// GetBlock returns a block by index.
func (bc *Blockchain) GetBlock(index uint64) (*Block, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if index >= uint64(len(bc.Chain)) {
		return nil, fmt.Errorf("block %d not found", index)
	}
	block := bc.Chain[index]
	return &block, nil
}

// GetLatestBlock returns the most recent block.
func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	block := bc.Chain[len(bc.Chain)-1]
	return &block
}

// Height returns the number of blocks in the chain.
func (bc *Blockchain) Height() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return uint64(len(bc.Chain))
}

// GetChain returns a copy of the full chain.
func (bc *Blockchain) GetChain() []Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	result := make([]Block, len(bc.Chain))
	copy(result, bc.Chain)
	return result
}

// CreditAddress adds balance to an address (used for faucet).
func (bc *Blockchain) CreditAddress(address string, amount uint64) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.Balances[address] += amount
	bc.store.SaveBalances(bc.Balances)
}

// TxResult wraps a transaction with the block it was found in.
type TxResult struct {
	Transaction Transaction `json:"transaction"`
	BlockIndex  uint64      `json:"blockIndex"`
	BlockHash   string      `json:"blockHash"`
}

// FindTransaction searches the chain for a transaction by ID.
func (bc *Blockchain) FindTransaction(txID string) (*TxResult, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.ID == txID {
				return &TxResult{
					Transaction: tx,
					BlockIndex:  block.Index,
					BlockHash:   block.Hash,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("transaction %s not found", txID)
}

// GetAddressTransactions returns all transactions involving an address.
func (bc *Blockchain) GetAddressTransactions(address string) []TxResult {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	var results []TxResult
	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.Sender == address || tx.Receiver == address {
				results = append(results, TxResult{
					Transaction: tx,
					BlockIndex:  block.Index,
					BlockHash:   block.Hash,
				})
			}
		}
	}
	return results
}
