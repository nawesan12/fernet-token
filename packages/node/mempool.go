package node

import (
	"sync"

	"github.com/nawesan12/fernet-token/packages/blockchain"
)

// Mempool is a thread-safe pending transaction pool.
type Mempool struct {
	mu   sync.RWMutex
	txns map[string]*blockchain.Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		txns: make(map[string]*blockchain.Transaction),
	}
}

// Add adds a transaction to the mempool.
func (m *Mempool) Add(tx *blockchain.Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.txns[tx.ID] = tx
}

// GetPending returns up to limit pending transactions.
func (m *Mempool) GetPending(limit int) []blockchain.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []blockchain.Transaction
	for _, tx := range m.txns {
		result = append(result, *tx)
		if len(result) >= limit {
			break
		}
	}
	return result
}

// RemoveConfirmed removes transactions that have been included in a block.
func (m *Mempool) RemoveConfirmed(txns []blockchain.Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, tx := range txns {
		delete(m.txns, tx.ID)
	}
}

// Count returns the number of pending transactions.
func (m *Mempool) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.txns)
}

// GetAll returns all pending transactions.
func (m *Mempool) GetAll() []blockchain.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []blockchain.Transaction
	for _, tx := range m.txns {
		result = append(result, *tx)
	}
	return result
}
