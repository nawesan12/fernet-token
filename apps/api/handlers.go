package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nawesan12/fernet-token/packages/blockchain"
	"github.com/nawesan12/fernet-token/packages/node"
	"github.com/nawesan12/fernet-token/packages/wallet"
)

type APIHandler struct {
	node *node.Node
}

func NewAPIHandler(n *node.Node) *APIHandler {
	return &APIHandler{node: n}
}

func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/blockchain", h.getBlockchain)
	mux.HandleFunc("GET /api/blockchain/height", h.getHeight)
	mux.HandleFunc("GET /api/block/{index}", h.getBlock)
	mux.HandleFunc("GET /api/balance/{address}", h.getBalance)
	mux.HandleFunc("GET /api/nonce/{address}", h.getNonce)
	mux.HandleFunc("GET /api/tx/pending", h.getPending)
	mux.HandleFunc("GET /api/tx/{id}", h.getTransaction)
	mux.HandleFunc("GET /api/address/{address}/transactions", h.getAddressTransactions)
	mux.HandleFunc("GET /api/peers", h.getPeers)
	mux.HandleFunc("POST /api/wallet/create", h.createWallet)
	mux.HandleFunc("POST /api/transaction", h.submitTransaction)
	mux.HandleFunc("POST /api/mine", h.mine)
	mux.HandleFunc("POST /api/peers/connect", h.connectPeer)
	mux.HandleFunc("POST /api/faucet", h.faucet)
}

func (h *APIHandler) getBlockchain(w http.ResponseWriter, r *http.Request) {
	chain := h.node.Blockchain.GetChain()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"chain":  chain,
		"height": len(chain),
	})
}

func (h *APIHandler) getHeight(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"height": h.node.Blockchain.Height(),
	})
}

func (h *APIHandler) getBlock(w http.ResponseWriter, r *http.Request) {
	indexStr := r.PathValue("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid block index")
		return
	}

	block, err := h.node.Blockchain.GetBlock(index)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, block)
}

func (h *APIHandler) getBalance(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	balance := h.node.Blockchain.GetBalance(address)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"address":  address,
		"balance":  balance,
		"formatted": formatFernet(balance),
	})
}

func (h *APIHandler) getNonce(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	nonce := h.node.Blockchain.GetNonce(address)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"address": address,
		"nonce":   nonce,
	})
}

func (h *APIHandler) getPending(w http.ResponseWriter, r *http.Request) {
	txns := h.node.Mempool.GetAll()
	if txns == nil {
		txns = []blockchain.Transaction{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"transactions": txns,
		"count":        len(txns),
	})
}

func (h *APIHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	txID := r.PathValue("id")
	result, err := h.node.Blockchain.FindTransaction(txID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *APIHandler) getAddressTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.PathValue("address")
	results := h.node.Blockchain.GetAddressTransactions(address)
	if results == nil {
		results = []blockchain.TxResult{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"address":      address,
		"transactions": results,
		"count":        len(results),
	})
}

func (h *APIHandler) getPeers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"peers": h.node.P2P.PeerAddresses(),
		"count": h.node.P2P.PeerCount(),
	})
}

func (h *APIHandler) createWallet(w http.ResponseWriter, r *http.Request) {
	wal, err := wallet.NewWallet()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create wallet")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"address":   wal.Address,
		"publicKey": wal.PublicKey,
	})
}

type submitTxRequest struct {
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

func (h *APIHandler) submitTransaction(w http.ResponseWriter, r *http.Request) {
	var req submitTxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tx := &blockchain.Transaction{
		ID:        req.ID,
		Sender:    req.Sender,
		Receiver:  req.Receiver,
		Amount:    req.Amount,
		Fee:       req.Fee,
		Nonce:     req.Nonce,
		Timestamp: req.Timestamp,
		PubKey:    req.PubKey,
		Signature: req.Signature,
	}

	if err := h.node.SubmitTransaction(tx); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]interface{}{
		"message": "transaction submitted",
		"txId":    tx.ID,
	})
}

type mineRequest struct {
	Miner string `json:"miner"`
}

func (h *APIHandler) mine(w http.ResponseWriter, r *http.Request) {
	var req mineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Miner == "" {
		writeError(w, http.StatusBadRequest, "miner address required")
		return
	}

	block, err := h.node.Mine(req.Miner)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "block mined",
		"block":   block,
	})
}

type connectPeerRequest struct {
	Address string `json:"address"`
}

func (h *APIHandler) connectPeer(w http.ResponseWriter, r *http.Request) {
	var req connectPeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Address == "" {
		writeError(w, http.StatusBadRequest, "peer address required")
		return
	}

	if err := h.node.P2P.ConnectToPeer(req.Address); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "connected to peer",
		"address": req.Address,
	})
}

type faucetRequest struct {
	Address string `json:"address"`
}

func (h *APIHandler) faucet(w http.ResponseWriter, r *http.Request) {
	var req faucetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Address == "" {
		writeError(w, http.StatusBadRequest, "address required")
		return
	}

	amount := uint64(100) * blockchain.OneFernet
	h.node.Blockchain.CreditAddress(req.Address, amount)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "faucet: credited 100 FERNET",
		"address": req.Address,
		"amount":  amount,
		"formatted": formatFernet(amount),
	})
}

func formatFernet(fernetoshi uint64) string {
	whole := fernetoshi / blockchain.OneFernet
	frac := fernetoshi % blockchain.OneFernet
	if frac == 0 {
		return fmt.Sprintf("%d FERNET", whole)
	}
	return fmt.Sprintf("%d.%08d FERNET", whole, frac)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
