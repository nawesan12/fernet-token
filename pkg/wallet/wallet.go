package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  string
	Address    string
}

// NewWallet generates a new ECDSA key pair and derives the wallet address.
func NewWallet() (*Wallet, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error generating wallet key: %w", err)
	}

	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	address := generateAddress(pubKey)

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  hex.EncodeToString(pubKey),
		Address:    address,
	}, nil
}

// generateAddress creates a wallet address using SHA-256.
func generateAddress(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	return hex.EncodeToString(hash[:20]) // Use first 20 bytes as address
}

// SignTransaction signs data using the wallet's private key.
func (w *Wallet) SignTransaction(data string) (string, string, error) {
	hash := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return "", "", fmt.Errorf("error signing transaction: %w", err)
	}
	return r.Text(16), s.Text(16), nil
}

// VerifySignature verifies if the signature is valid for the given data.
func VerifySignature(pubKeyHex, data, rHex, sHex string) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubKeyBytes) != 64 {
		return false, errors.New("invalid public key")
	}

	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])
	pubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	hash := sha256.Sum256([]byte(data))
	r := new(big.Int)
	s := new(big.Int)
	r.SetString(rHex, 16)
	s.SetString(sHex, 16)

	return ecdsa.Verify(pubKey, hash[:], r, s), nil
}
