package blockchain

// 1 FERNET = 100,000,000 fernetoshi
const (
	Fernetoshi     uint64 = 1
	OneFernet      uint64 = 100_000_000
	TotalSupply    uint64 = 21_000_000 * OneFernet
	MiningReward   uint64 = 50 * OneFernet
	MaxTxPerBlock  int    = 100
	TargetPrefix          = "0000"
	CoinbaseSender        = "0000000000000000000000000000000000000000"

	// Fixed genesis timestamp for deterministic genesis block
	GenesisTimestamp int64 = 1700000000
)

// Block represents a block in the blockchain.
type Block struct {
	Index        uint64        `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prevHash"`
	Hash         string        `json:"hash"`
	Nonce        uint64        `json:"nonce"`
	Miner        string        `json:"miner"`
}

// BlockHashData is the same as Block but without the Hash field,
// used to avoid circular hash calculation.
type BlockHashData struct {
	Index        uint64        `json:"index"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prevHash"`
	Nonce        uint64        `json:"nonce"`
	Miner        string        `json:"miner"`
}

// Transaction represents a transfer of fernetoshi between addresses.
type Transaction struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Amount    uint64 `json:"amount"`
	Fee       uint64 `json:"fee"`
	Nonce     uint64 `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	PubKey    string `json:"pubKey"`
	Signature string `json:"signature"`
}

// Storage is the persistence interface for the blockchain.
type Storage interface {
	SaveBlock(block Block) error
	LoadChain() ([]Block, error)
	SaveBalances(balances map[string]uint64) error
	LoadBalances() (map[string]uint64, error)
	SaveNonces(nonces map[string]uint64) error
	LoadNonces() (map[string]uint64, error)
	Close() error
}
