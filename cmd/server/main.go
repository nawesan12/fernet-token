// cmd/server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nawesan12/olive-token/pkg/blockchain"
	"github.com/nawesan12/olive-token/pkg/p2p"
	"github.com/nawesan12/olive-token/pkg/wallet"
)

var (
	bc     = blockchain.NewBlockchain()
	server = p2p.NewP2PServer()
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [port]")
	}
	port := os.Args[1]

	// Start P2P server
	go server.Start(port)

	// Setup HTTP routes
	http.HandleFunc("/wallet/create", createWalletHandler)
	http.HandleFunc("/wallet/balance", getBalanceHandler)
	http.HandleFunc("/transaction/send", sendTransactionHandler)
	http.HandleFunc("/blockchain", getBlockchainHandler)
	http.HandleFunc("/blockchain/latest", getLatestBlockHandler)
	http.HandleFunc("/peer/connect", connectPeerHandler)
	http.HandleFunc("/peer/list", listPeersHandler)

	log.Println("🚀 Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createWalletHandler(w http.ResponseWriter, r *http.Request) {
	wallet := wallet.NewWallet()
	if wallet == nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(wallet.PublicKey)
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for wallet balance retrieval
	json.NewEncoder(w).Encode("Balance feature not implemented yet")
}

func sendTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var tx blockchain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if blockchain.SimpleConsensus(&tx) {
		bc.AddTransaction(&tx)
		json.NewEncoder(w).Encode("✅ Transaction added to blockchain")
	} else {
		http.Error(w, "❌ Invalid transaction", http.StatusBadRequest)
	}
}

func getBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc)
}

func getLatestBlockHandler(w http.ResponseWriter, r *http.Request) {
	if len(bc.Transactions) == 0 {
		http.Error(w, "Blockchain is empty", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(bc.Transactions[len(bc.Transactions)-1])
}

func connectPeerHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Peer connection not implemented yet")
}

func listPeersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(server.Peers)
}
