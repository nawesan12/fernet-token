// cmd/server/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/nawesan12/fernet-token/pkg/blockchain"
	"github.com/nawesan12/fernet-token/pkg/p2p"
	"github.com/nawesan12/fernet-token/pkg/wallet"
)

var (
	bc           = blockchain.NewBlockchain()
	server       = p2p.NewP2PServer()
	balances     = make(map[string]float64)
	totalSupply  = 1000000.0
	miningReward = 10.0
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [port]")
	}
	port := os.Args[1]

	// Distribute initial supply
	genesisWallet := wallet.NewWallet()
	balances[genesisWallet.PublicKey] = totalSupply
	fmt.Println("Genesis wallet created with", totalSupply, "tokens")

	// Start P2P server
	go server.Start(port)

	// Setup HTTP routes
	http.HandleFunc("/wallet/create", createWalletHandler)
	http.HandleFunc("/wallet/balance", getBalanceHandler)
	http.HandleFunc("/transaction/send", sendTransactionHandler)
	http.HandleFunc("/blockchain", getBlockchainHandler)
	http.HandleFunc("/peer/connect", connectPeerHandler)
	http.HandleFunc("/peer/list", listPeersHandler)
	http.HandleFunc("/mine", mineBlockHandler)
	http.HandleFunc("/transaction/add", handleAddTransaction)
	http.HandleFunc("/transaction/pending", handleGetPendingTransactions)

	server := &http.Server{Addr: ":" + port}
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
	log.Println("🚀 Server running on port", ln.Addr().(*net.TCPAddr).Port)
	log.Fatal(server.Serve(ln))
}

func createWalletHandler(w http.ResponseWriter, r *http.Request) {
	wallet := wallet.NewWallet()
	if wallet == nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}
	balances[wallet.PublicKey] = 0
	json.NewEncoder(w).Encode(wallet.PublicKey)
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for wallet balance retrieval
	addr := r.URL.Query().Get("address")
	balance, exists := balances[addr]
	if !exists {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func sendTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var tx TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if balances[tx.Sender] < tx.Amount {
		http.Error(w, "Insufficient balance", http.StatusForbidden)
		return
	}
	balances[tx.Sender] -= tx.Amount
	balances[tx.Receiver] += tx.Amount
	bc.AddTransaction(tx.Sender, tx.Receiver, tx.Amount)
	json.NewEncoder(w).Encode("Transaction processed successfully")
}

func getBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc)
}

func connectPeerHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Peer connection not implemented yet")
}

func listPeersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(server.Peers)
}

func mineBlockHandler(w http.ResponseWriter, r *http.Request) {
	minerAddress := r.URL.Query().Get("miner")
	if minerAddress == "" {
		http.Error(w, "Miner address not provided", http.StatusBadRequest)
		return
	}

	bc.MineBlock(minerAddress, miningReward)
	balances[minerAddress] += miningReward
	json.NewEncoder(w).Encode("Block mined and reward issued to miner")
}

type TransactionRequest struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

func handleAddTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	bc.AddTransaction(req.Sender, req.Receiver, req.Amount)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction added"})
}

func handleGetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc.GetPendingTransactions())
}
