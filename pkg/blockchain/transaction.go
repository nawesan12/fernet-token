package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Transaction struct {
	ID        string  // Hash de la transacción
	Sender    string  // Dirección del remitente
	Receiver  string  // Dirección del receptor
	Amount    float64 // Monto transferido
	Timestamp int64   // Marca de tiempo
}

func NewTransaction(sender, receiver string, amount float64) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
	tx.ID = tx.CalculateHash()
	return tx
}

func (t *Transaction) CalculateHash() string {
	data := fmt.Sprintf("%s%s%f%d", t.Sender, t.Receiver, t.Amount, t.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
