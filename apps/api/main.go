package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nawesan12/fernet-token/packages/blockchain"
	"github.com/nawesan12/fernet-token/packages/p2p"
	"github.com/nawesan12/fernet-token/packages/wallet"
)

var (
	bc           *blockchain.Blockchain
	server       *p2p.P2PServer
	balances     = make(map[string]float64)
	totalSupply  = 1000000.0
	miningReward = 10.0
	mutex        sync.RWMutex
)

func main() {
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	initializeBlockchain()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleShutdown(cancel)
	go server.Start(*port)

	startHTTPServer(*port, ctx)
}

func initializeBlockchain() {
	bc = blockchain.NewBlockchain()
	server = p2p.NewP2PServer()

	genesisWallet, err := wallet.NewWallet()
	if err != nil {
		log.Fatal("Failed to create genesis wallet")
	}

	balances[genesisWallet.PublicKey] = totalSupply
	log.Println("✅ Genesis wallet created with", totalSupply, "$Fernet")
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Shutting down gracefully...")
	cancel()
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

func startHTTPServer(port string, ctx context.Context) {
	mux := http.NewServeMux()
	initializeRoutes(mux)

	server := &http.Server{Addr: ":" + port, Handler: mux}
	go func() {
		log.Println("🚀 HTTP Server running on port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("❌ Server error:", err)
		}
	}()

	<-ctx.Done()
	shutdownServer(server)
}

func shutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("❌ Error shutting down server:", err)
	}
	log.Println("✅ Server shut down gracefully")
}

func initializeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/wallet/create", createWalletHandler)
	mux.HandleFunc("/wallet/balance", getBalanceHandler)
	mux.HandleFunc("/transaction/send", sendTransactionHandler)
	mux.HandleFunc("/blockchain", getBlockchainHandler)
	mux.HandleFunc("/peer/connect", connectPeerHandler)
	mux.HandleFunc("/peer/list", listPeersHandler)
	mux.HandleFunc("/mine", mineBlockHandler)
	mux.HandleFunc("/transaction/add", handleAddTransaction)
	mux.HandleFunc("/transaction/pending", handleGetPendingTransactions)
}

func createWalletHandler(w http.ResponseWriter, r *http.Request) {
	wallet, err := wallet.NewWallet()
	if err != nil {
		http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	balances[wallet.PublicKey] = 0
	mutex.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"publicKey": wallet.PublicKey})
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	mutex.RLock()
	balance, exists := balances[addr]
	mutex.RUnlock()

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

	mutex.Lock()
	defer mutex.Unlock()

	if balances[tx.Sender] < tx.Amount {
		http.Error(w, "Insufficient balance", http.StatusForbidden)
		return
	}

	balances[tx.Sender] -= tx.Amount
	balances[tx.Receiver] += tx.Amount
	bc.AddTransaction(tx.Sender, tx.Receiver, tx.Amount)

	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction processed successfully"})
}

func getBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc)
}

func connectPeerHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Peer connection not implemented yet", http.StatusNotImplemented)
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

	mutex.Lock()
	bc.MineBlock(minerAddress, miningReward)
	balances[minerAddress] += miningReward
	mutex.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"message": "Block mined and reward issued to miner"})
}

func handleAddTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	bc.AddTransaction(req.Sender, req.Receiver, req.Amount)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction added"})
}

func handleGetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(bc.GetPendingTransactions())
}

type TransactionRequest struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}
