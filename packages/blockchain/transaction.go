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
	Timestamp int64   // Marca de tiempo en nanosegundos
}

// NewTransaction crea una nueva transacción con un hash único
func NewTransaction(sender, receiver string, amount float64) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().UTC().UnixNano(),
	}
	tx.ID = tx.CalculateHash()
	return tx
}

// CalculateHash genera un hash SHA-256 único para la transacción
func (t *Transaction) CalculateHash() string {
	data := fmt.Sprintf("%s%s%.8f%d", t.Sender, t.Receiver, t.Amount, t.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// IsValid verifica si la transacción es válida
func (t *Transaction) IsValid() bool {
	return t.Sender != "" && t.Receiver != "" && t.Amount > 0
}
