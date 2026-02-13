package blockchain

import (
	"fmt"
	"strings"
)

// ValidateBlock validates a block received from a peer.
func (bc *Blockchain) ValidateBlock(block *Block) error {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	latestBlock := bc.Chain[len(bc.Chain)-1]

	// Check index sequence
	if block.Index != latestBlock.Index+1 {
		return fmt.Errorf("expected index %d, got %d", latestBlock.Index+1, block.Index)
	}

	// Check prev hash link
	if block.PrevHash != latestBlock.Hash {
		return fmt.Errorf("prev hash mismatch: expected %s, got %s", latestBlock.Hash, block.PrevHash)
	}

	// Check hash correctness
	expectedHash := CalculateBlockHash(block)
	if block.Hash != expectedHash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, block.Hash)
	}

	// Check PoW
	if !strings.HasPrefix(block.Hash, TargetPrefix) {
		return fmt.Errorf("insufficient proof of work")
	}

	// Check coinbase
	if len(block.Transactions) == 0 || block.Transactions[0].Sender != CoinbaseSender {
		return fmt.Errorf("missing coinbase transaction")
	}
	if block.Transactions[0].Amount != MiningReward {
		return fmt.Errorf("invalid coinbase reward: expected %d, got %d", MiningReward, block.Transactions[0].Amount)
	}

	// Verify all non-coinbase transaction signatures
	for i := 1; i < len(block.Transactions); i++ {
		if err := block.Transactions[i].VerifySignature(); err != nil {
			return fmt.Errorf("tx %d signature invalid: %w", i, err)
		}
	}

	return nil
}

// AddBlock validates then appends a block received from a peer.
func (bc *Blockchain) AddBlock(block *Block) error {
	if err := bc.ValidateBlock(block); err != nil {
		return fmt.Errorf("block validation failed: %w", err)
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Apply state changes
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

	bc.Chain = append(bc.Chain, *block)
	bc.store.SaveBlock(*block)
	bc.store.SaveBalances(bc.Balances)
	bc.store.SaveNonces(bc.Nonces)

	return nil
}

// ShouldReplaceChain returns true if the given chain is longer and valid.
func (bc *Blockchain) ShouldReplaceChain(newChain []Block) bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(newChain) <= len(bc.Chain) {
		return false
	}

	// Validate the new chain
	if len(newChain) == 0 || newChain[0].PrevHash != "0" {
		return false
	}

	for i := 1; i < len(newChain); i++ {
		block := newChain[i]
		prevBlock := newChain[i-1]

		if block.Index != prevBlock.Index+1 {
			return false
		}
		if block.PrevHash != prevBlock.Hash {
			return false
		}
		expectedHash := CalculateBlockHash(&block)
		if block.Hash != expectedHash {
			return false
		}
		if !strings.HasPrefix(block.Hash, TargetPrefix) && i > 0 {
			return false
		}
	}

	return true
}

// ReplaceChain replaces the current chain with a longer valid chain.
func (bc *Blockchain) ReplaceChain(newChain []Block) error {
	if !bc.ShouldReplaceChain(newChain) {
		return fmt.Errorf("new chain is not valid or not longer")
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	bc.Chain = newChain

	// Rebuild state from scratch
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
		bc.store.SaveBlock(block)
	}

	bc.store.SaveBalances(bc.Balances)
	bc.store.SaveNonces(bc.Nonces)

	return nil
}
