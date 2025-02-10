package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
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

// NewTransaction crea una nueva transacción
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

// CalculateHash genera el hash de la transacción
func (t *Transaction) CalculateHash() string {
	// Convert Amount to string
	amountStr := fmt.Sprintf("%f", t.Amount)
	data := t.Sender + t.Receiver + amountStr + string(t.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
