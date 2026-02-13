package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
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

	pubKey := marshalPublicKey(&privKey.PublicKey)
	address := generateAddress(pubKey)

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  hex.EncodeToString(pubKey),
		Address:    address,
	}, nil
}

// marshalPublicKey returns X||Y as exactly 64 bytes (zero-padded).
func marshalPublicKey(pub *ecdsa.PublicKey) []byte {
	x := zeroPad(pub.X.Bytes(), 32)
	y := zeroPad(pub.Y.Bytes(), 32)
	return append(x, y...)
}

// zeroPad pads a byte slice to the given length with leading zeros.
func zeroPad(b []byte, length int) []byte {
	if len(b) >= length {
		return b[:length]
	}
	padded := make([]byte, length)
	copy(padded[length-len(b):], b)
	return padded
}

// generateAddress creates a wallet address using SHA-256 (first 20 bytes).
func generateAddress(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	return hex.EncodeToString(hash[:20])
}

// Sign signs data bytes and returns a fixed-size R||S hex string (128 chars).
func (w *Wallet) Sign(data []byte) (string, error) {
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, data)
	if err != nil {
		return "", fmt.Errorf("error signing: %w", err)
	}

	rBytes := zeroPad(r.Bytes(), 32)
	sBytes := zeroPad(s.Bytes(), 32)
	sig := append(rBytes, sBytes...)
	return hex.EncodeToString(sig), nil
}

// VerifySignature verifies an ECDSA signature given a public key hex, data, and signature hex.
func VerifySignature(pubKeyHex string, data []byte, sigHex string) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubKeyBytes) != 64 {
		return false, errors.New("invalid public key")
	}

	x := new(big.Int).SetBytes(pubKeyBytes[:32])
	y := new(big.Int).SetBytes(pubKeyBytes[32:])
	pubKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil || len(sigBytes) != 64 {
		return false, errors.New("invalid signature")
	}

	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	return ecdsa.Verify(pubKey, data, r, s), nil
}

// SaveToFile saves the wallet's private key to a PEM file.
func (w *Wallet) SaveToFile(path string) error {
	der, err := x509.MarshalECPrivateKey(w.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}

	return os.WriteFile(path, pem.EncodeToMemory(block), 0600)
}

// LoadFromFile loads a wallet from a PEM private key file.
func LoadFromFile(path string) (*Wallet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}

	pubKey := marshalPublicKey(&privKey.PublicKey)
	address := generateAddress(pubKey)

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  hex.EncodeToString(pubKey),
		Address:    address,
	}, nil
}
