package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// NewTransaction creates a new unsigned transaction.
func NewTransaction(sender, receiver string, amount, fee, nonce uint64, pubKey string) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Fee:       fee,
		Nonce:     nonce,
		Timestamp: time.Now().UTC().UnixNano(),
		PubKey:    pubKey,
	}
	tx.ID = tx.CalculateHash()
	return tx
}

// NewCoinbaseTx creates a coinbase (mining reward) transaction.
func NewCoinbaseTx(receiver string, reward uint64) *Transaction {
	tx := &Transaction{
		Sender:    CoinbaseSender,
		Receiver:  receiver,
		Amount:    reward,
		Fee:       0,
		Nonce:     0,
		Timestamp: time.Now().UTC().UnixNano(),
		PubKey:    "",
		Signature: "",
	}
	tx.ID = tx.CalculateHash()
	return tx
}

// CalculateHash computes the SHA-256 hash of the transaction's core data.
func (t *Transaction) CalculateHash() string {
	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d", t.Sender, t.Receiver, t.Amount, t.Fee, t.Nonce, t.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SignableData returns the SHA-256 hash bytes used for signing.
func (t *Transaction) SignableData() []byte {
	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d", t.Sender, t.Receiver, t.Amount, t.Fee, t.Nonce, t.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// VerifySignature verifies the transaction's ECDSA signature.
func (t *Transaction) VerifySignature() error {
	if t.Sender == CoinbaseSender {
		return nil
	}

	if t.PubKey == "" || t.Signature == "" {
		return errors.New("missing public key or signature")
	}

	pubKeyBytes, err := hex.DecodeString(t.PubKey)
	if err != nil || len(pubKeyBytes) != 64 {
		return errors.New("invalid public key format")
	}

	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])
	pubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	// Verify public key derives to sender address
	addrHash := sha256.Sum256(pubKeyBytes)
	derivedAddr := hex.EncodeToString(addrHash[:20])
	if derivedAddr != t.Sender {
		return errors.New("public key does not match sender address")
	}

	// Decode signature (R||S, 64 bytes = 128 hex chars)
	sigBytes, err := hex.DecodeString(t.Signature)
	if err != nil || len(sigBytes) != 64 {
		return fmt.Errorf("invalid signature format: expected 128 hex chars, got %d", len(t.Signature))
	}

	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	hash := t.SignableData()
	if !ecdsa.Verify(pubKey, hash, r, s) {
		return errors.New("signature verification failed")
	}

	return nil
}

// IsValid checks all transaction validity rules.
func (t *Transaction) IsValid() error {
	if t.Sender == "" {
		return errors.New("sender is empty")
	}
	if t.Receiver == "" {
		return errors.New("receiver is empty")
	}
	if t.Sender == t.Receiver && t.Sender != CoinbaseSender {
		return errors.New("sender and receiver are the same")
	}
	if t.Amount == 0 {
		return errors.New("amount is zero")
	}

	expectedHash := t.CalculateHash()
	if t.ID != expectedHash {
		return errors.New("transaction hash mismatch")
	}

	if err := t.VerifySignature(); err != nil {
		return fmt.Errorf("signature error: %w", err)
	}

	return nil
}
