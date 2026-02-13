package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func generateTestWallet() (*ecdsa.PrivateKey, string, string) {
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	x := zeroPadTest(privKey.PublicKey.X.Bytes(), 32)
	y := zeroPadTest(privKey.PublicKey.Y.Bytes(), 32)
	pubKeyBytes := append(x, y...)
	pubKeyHex := hex.EncodeToString(pubKeyBytes)

	addrHash := sha256.Sum256(pubKeyBytes)
	address := hex.EncodeToString(addrHash[:20])

	return privKey, pubKeyHex, address
}

func zeroPadTest(b []byte, length int) []byte {
	if len(b) >= length {
		return b[:length]
	}
	padded := make([]byte, length)
	copy(padded[length-len(b):], b)
	return padded
}

func signTx(privKey *ecdsa.PrivateKey, tx *Transaction) {
	hash := tx.SignableData()
	r, s, _ := ecdsa.Sign(rand.Reader, privKey, hash)

	rBytes := zeroPadTest(r.Bytes(), 32)
	sBytes := zeroPadTest(s.Bytes(), 32)
	sig := append(rBytes, sBytes...)
	tx.Signature = hex.EncodeToString(sig)
}

func TestSignAndVerifyTransaction(t *testing.T) {
	privKey, pubKeyHex, senderAddr := generateTestWallet()
	_, _, receiverAddr := generateTestWallet()

	tx := NewTransaction(senderAddr, receiverAddr, 100*OneFernet, OneFernet, 0, pubKeyHex)
	signTx(privKey, tx)
	// Recalculate ID after all fields are set
	tx.ID = tx.CalculateHash()

	if err := tx.VerifySignature(); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}

	if err := tx.IsValid(); err != nil {
		t.Fatalf("transaction should be valid: %v", err)
	}
}

func TestCoinbaseTransaction(t *testing.T) {
	tx := NewCoinbaseTx("someaddress", MiningReward)

	if tx.Sender != CoinbaseSender {
		t.Errorf("coinbase sender should be %s, got %s", CoinbaseSender, tx.Sender)
	}

	if err := tx.VerifySignature(); err != nil {
		t.Errorf("coinbase should pass signature check: %v", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	_, pubKeyHex, senderAddr := generateTestWallet()
	_, _, receiverAddr := generateTestWallet()

	tx := NewTransaction(senderAddr, receiverAddr, 100*OneFernet, OneFernet, 0, pubKeyHex)
	tx.Signature = "0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000"
	tx.ID = tx.CalculateHash()

	if err := tx.VerifySignature(); err == nil {
		t.Error("invalid signature should fail verification")
	}
}

func TestWrongSenderAddress(t *testing.T) {
	privKey, pubKeyHex, _ := generateTestWallet()
	_, _, receiverAddr := generateTestWallet()

	// Use wrong sender address
	tx := NewTransaction("wrongaddress1234567890", receiverAddr, 100*OneFernet, OneFernet, 0, pubKeyHex)
	signTx(privKey, tx)
	tx.ID = tx.CalculateHash()

	if err := tx.VerifySignature(); err == nil {
		t.Error("wrong sender should fail verification")
	}
}

func TestTransactionHash(t *testing.T) {
	tx := NewTransaction("sender", "receiver", 100, 1, 0, "pubkey")
	hash := tx.CalculateHash()

	if hash != tx.ID {
		t.Error("ID should match CalculateHash result")
	}

	// Same params should produce same hash
	tx2 := &Transaction{
		Sender:    tx.Sender,
		Receiver:  tx.Receiver,
		Amount:    tx.Amount,
		Fee:       tx.Fee,
		Nonce:     tx.Nonce,
		Timestamp: tx.Timestamp,
		PubKey:    tx.PubKey,
	}
	if tx2.CalculateHash() != hash {
		t.Error("identical transaction data should produce identical hash")
	}
}
