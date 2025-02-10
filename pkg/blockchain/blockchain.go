package blockchain

import "sync"

type Blockchain struct {
	Transactions []*Transaction
	Mutex        sync.Mutex
}

// NewBlockchain crea una nueva instancia de Blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Transactions: []*Transaction{},
	}
}

// AddTransaction agrega una transacción a la Blockchain
func (bc *Blockchain) AddTransaction(tx *Transaction) {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()
	bc.Transactions = append(bc.Transactions, tx)
}
