package blockchain

import (
	"strings"
	"testing"
)

func TestGenesisBlock(t *testing.T) {
	store := NewMemoryStorage()
	bc, err := NewBlockchain(store)
	if err != nil {
		t.Fatalf("NewBlockchain failed: %v", err)
	}

	if len(bc.Chain) != 1 {
		t.Fatalf("expected 1 block, got %d", len(bc.Chain))
	}

	genesis := bc.Chain[0]
	if genesis.Index != 0 {
		t.Errorf("genesis index should be 0, got %d", genesis.Index)
	}
	if genesis.PrevHash != "0" {
		t.Errorf("genesis prevHash should be '0', got '%s'", genesis.PrevHash)
	}
	if genesis.Timestamp != GenesisTimestamp {
		t.Errorf("genesis timestamp should be %d, got %d", GenesisTimestamp, genesis.Timestamp)
	}
}

func TestDeterministicGenesis(t *testing.T) {
	store1 := NewMemoryStorage()
	bc1, _ := NewBlockchain(store1)

	store2 := NewMemoryStorage()
	bc2, _ := NewBlockchain(store2)

	if bc1.Chain[0].Hash != bc2.Chain[0].Hash {
		t.Errorf("genesis hashes should be identical: %s != %s", bc1.Chain[0].Hash, bc2.Chain[0].Hash)
	}
}

func TestCalculateBlockHash(t *testing.T) {
	block := &Block{
		Index:        1,
		Timestamp:    1000000,
		Transactions: []Transaction{},
		PrevHash:     "abc",
		Hash:         "should-be-excluded",
		Nonce:        42,
		Miner:        "miner1",
	}

	hash1 := CalculateBlockHash(block)
	block.Hash = "different-value"
	hash2 := CalculateBlockHash(block)

	if hash1 != hash2 {
		t.Error("hash should not depend on the Hash field (circular hash bug)")
	}
}

func TestMineBlock(t *testing.T) {
	store := NewMemoryStorage()
	bc, _ := NewBlockchain(store)

	block, err := bc.MineBlock("miner1", nil)
	if err != nil {
		t.Fatalf("MineBlock failed: %v", err)
	}

	if block.Index != 1 {
		t.Errorf("expected block index 1, got %d", block.Index)
	}
	if !strings.HasPrefix(block.Hash, TargetPrefix) {
		t.Errorf("block hash should start with %s, got %s", TargetPrefix, block.Hash)
	}
	if block.Miner != "miner1" {
		t.Errorf("miner should be 'miner1', got '%s'", block.Miner)
	}

	// Check coinbase
	if len(block.Transactions) < 1 {
		t.Fatal("block should have at least coinbase tx")
	}
	if block.Transactions[0].Sender != CoinbaseSender {
		t.Error("first tx should be coinbase")
	}
	if block.Transactions[0].Amount != MiningReward {
		t.Errorf("coinbase amount should be %d, got %d", MiningReward, block.Transactions[0].Amount)
	}

	// Check miner balance
	if bc.GetBalance("miner1") != MiningReward {
		t.Errorf("miner balance should be %d, got %d", MiningReward, bc.GetBalance("miner1"))
	}
}

func TestPersistence(t *testing.T) {
	store := NewMemoryStorage()
	bc1, _ := NewBlockchain(store)
	bc1.MineBlock("miner1", nil)

	// Load from same store
	bc2, _ := NewBlockchain(store)
	if len(bc2.Chain) != 2 {
		t.Errorf("expected 2 blocks after reload, got %d", len(bc2.Chain))
	}
	if bc2.GetBalance("miner1") != MiningReward {
		t.Errorf("balance should persist: expected %d, got %d", MiningReward, bc2.GetBalance("miner1"))
	}
}

func TestValidateChain(t *testing.T) {
	store := NewMemoryStorage()
	bc, _ := NewBlockchain(store)
	bc.MineBlock("miner1", nil)
	bc.MineBlock("miner1", nil)

	if err := bc.ValidateChain(); err != nil {
		t.Errorf("valid chain should pass validation: %v", err)
	}
}

func TestChainHeight(t *testing.T) {
	store := NewMemoryStorage()
	bc, _ := NewBlockchain(store)

	if bc.Height() != 1 {
		t.Errorf("initial height should be 1 (genesis), got %d", bc.Height())
	}

	bc.MineBlock("miner1", nil)
	if bc.Height() != 2 {
		t.Errorf("height after mining should be 2, got %d", bc.Height())
	}
}
