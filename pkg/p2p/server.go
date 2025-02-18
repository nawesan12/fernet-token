package p2p

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	Address string
	Conn    net.Conn
}

type P2PServer struct {
	Peers map[string]*Peer
	Mutex sync.Mutex
}

func NewP2PServer() *P2PServer {
	return &P2PServer{Peers: make(map[string]*Peer)}
}

func (s *P2PServer) Start(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("❌ Error starting P2P server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("🔗 P2P Server listening on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("⚠️ Connection error:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *P2PServer) handleConnection(conn net.Conn) {
	peerAddr := conn.RemoteAddr().String()
	s.Mutex.Lock()
	s.Peers[peerAddr] = &Peer{Address: peerAddr, Conn: conn}
	s.Mutex.Unlock()

	fmt.Println("✅ New peer connected:", peerAddr)
	defer conn.Close()

	// Listen for messages
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("🔌 Peer disconnected:", peerAddr)
			s.Mutex.Lock()
			delete(s.Peers, peerAddr)
			s.Mutex.Unlock()
			return
		}
		s.handleMessage(buffer[:n])
	}
}

func (s *P2PServer) handleMessage(data []byte) {
	var message map[string]interface{}
	if err := json.Unmarshal(data, &message); err != nil {
		fmt.Println("❌ Failed to parse message:", err)
		return
	}
	fmt.Println("📩 Received message:", message)
	// Handle different message types here (e.g., transactions, blocks, peer requests)
}

func (s *P2PServer) BroadcastMessage(msg map[string]interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("❌ Failed to serialize message:", err)
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	for _, peer := range s.Peers {
		_, err := peer.Conn.Write(data)
		if err != nil {
			fmt.Println("⚠️ Failed to send message to", peer.Address, err)
		}
	}
}
