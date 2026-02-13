package blockchain

import (
	"encoding/json"
	"fmt"
	"sync"

	bolt "go.etcd.io/bbolt"
)

var (
	blocksBucket   = []byte("blocks")
	balancesBucket = []byte("balances")
	noncesBucket   = []byte("nonces")
	metaBucket     = []byte("meta")
)

// BoltStorage implements Storage using bbolt (single-file embedded DB).
type BoltStorage struct {
	db *bolt.DB
}

func NewBoltStorage(path string) (*BoltStorage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range [][]byte{blocksBucket, balancesBucket, noncesBucket, metaBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create buckets: %w", err)
	}

	return &BoltStorage{db: db}, nil
}

func (s *BoltStorage) SaveBlock(block Block) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		data, err := json.Marshal(block)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%010d", block.Index)
		return b.Put([]byte(key), data)
	})
}

func (s *BoltStorage) LoadChain() ([]Block, error) {
	var chain []Block
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		return b.ForEach(func(k, v []byte) error {
			var block Block
			if err := json.Unmarshal(v, &block); err != nil {
				return err
			}
			chain = append(chain, block)
			return nil
		})
	})
	return chain, err
}

func (s *BoltStorage) SaveBalances(balances map[string]uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(balancesBucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		for addr, bal := range balances {
			data, err := json.Marshal(bal)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(addr), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BoltStorage) LoadBalances() (map[string]uint64, error) {
	balances := make(map[string]uint64)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(balancesBucket)
		return b.ForEach(func(k, v []byte) error {
			var bal uint64
			if err := json.Unmarshal(v, &bal); err != nil {
				return err
			}
			balances[string(k)] = bal
			return nil
		})
	})
	return balances, err
}

func (s *BoltStorage) SaveNonces(nonces map[string]uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(noncesBucket)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		for addr, nonce := range nonces {
			data, err := json.Marshal(nonce)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(addr), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BoltStorage) LoadNonces() (map[string]uint64, error) {
	nonces := make(map[string]uint64)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(noncesBucket)
		return b.ForEach(func(k, v []byte) error {
			var nonce uint64
			if err := json.Unmarshal(v, &nonce); err != nil {
				return err
			}
			nonces[string(k)] = nonce
			return nil
		})
	})
	return nonces, err
}

func (s *BoltStorage) Close() error {
	return s.db.Close()
}

// MemoryStorage implements Storage in-memory for testing.
type MemoryStorage struct {
	mu       sync.Mutex
	blocks   []Block
	balances map[string]uint64
	nonces   map[string]uint64
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		balances: make(map[string]uint64),
		nonces:   make(map[string]uint64),
	}
}

func (s *MemoryStorage) SaveBlock(block Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blocks = append(s.blocks, block)
	return nil
}

func (s *MemoryStorage) LoadChain() ([]Block, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Block, len(s.blocks))
	copy(result, s.blocks)
	return result, nil
}

func (s *MemoryStorage) SaveBalances(balances map[string]uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.balances = make(map[string]uint64)
	for k, v := range balances {
		s.balances[k] = v
	}
	return nil
}

func (s *MemoryStorage) LoadBalances() (map[string]uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]uint64)
	for k, v := range s.balances {
		result[k] = v
	}
	return result, nil
}

func (s *MemoryStorage) SaveNonces(nonces map[string]uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nonces = make(map[string]uint64)
	for k, v := range nonces {
		s.nonces[k] = v
	}
	return nil
}

func (s *MemoryStorage) LoadNonces() (map[string]uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make(map[string]uint64)
	for k, v := range s.nonces {
		result[k] = v
	}
	return result, nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
