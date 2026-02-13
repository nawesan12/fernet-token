package p2p

import (
	"log"
	"net"
	"sync"

	"github.com/nawesan12/fernet-token/packages/blockchain"
)

// MessageHandler is called when a message is received from a peer.
type MessageHandler func(Message)

type peer struct {
	addr string
	conn net.Conn
}

type P2PServer struct {
	port     string
	handler  MessageHandler
	peers    map[string]*peer
	mu       sync.RWMutex
	listener net.Listener
	quit     chan struct{}
}

func NewP2PServer(port string, handler MessageHandler) *P2PServer {
	return &P2PServer{
		port:    port,
		handler: handler,
		peers:   make(map[string]*peer),
		quit:    make(chan struct{}),
	}
}

// Start listens for incoming TCP connections.
func (s *P2PServer) Start() {
	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Printf("P2P: failed to listen on port %s: %v", s.port, err)
		return
	}
	s.listener = ln
	log.Printf("P2P: listening on port %s", s.port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Printf("P2P: accept error: %v", err)
				continue
			}
		}
		go s.handleConn(conn)
	}
}

func (s *P2PServer) handleConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	s.mu.Lock()
	s.peers[addr] = &peer{addr: addr, conn: conn}
	s.mu.Unlock()

	log.Printf("P2P: peer connected: %s", addr)

	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.peers, addr)
		s.mu.Unlock()
		log.Printf("P2P: peer disconnected: %s", addr)
	}()

	for {
		msg, err := ReadMessage(conn)
		if err != nil {
			return
		}

		if msg.Type == MsgPing {
			WriteMessage(conn, Message{Type: MsgPong})
			continue
		}

		msg.SenderAddr = addr
		s.handler(msg)
	}
}

// ConnectToPeer establishes a persistent outbound connection to a peer.
func (s *P2PServer) ConnectToPeer(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.peers[address] = &peer{addr: address, conn: conn}
	s.mu.Unlock()

	log.Printf("P2P: connected to peer: %s", address)

	// Start listening for messages from this peer
	go func() {
		defer func() {
			conn.Close()
			s.mu.Lock()
			delete(s.peers, address)
			s.mu.Unlock()
			log.Printf("P2P: outbound peer disconnected: %s", address)
		}()

		for {
			msg, err := ReadMessage(conn)
			if err != nil {
				return
			}

			if msg.Type == MsgPing {
				WriteMessage(conn, Message{Type: MsgPong})
				continue
			}

			msg.SenderAddr = address
			s.handler(msg)
		}
	}()

	// Request blocks from the peer
	WriteMessage(conn, Message{Type: MsgGetBlocks})

	return nil
}

// BroadcastTransaction sends a transaction to all connected peers.
func (s *P2PServer) BroadcastTransaction(tx *blockchain.Transaction) {
	msg := Message{Type: MsgTransaction, Transaction: tx}
	s.broadcast(msg)
}

// BroadcastBlock sends a block to all connected peers.
func (s *P2PServer) BroadcastBlock(block *blockchain.Block) {
	msg := Message{Type: MsgBlock, Block: block}
	s.broadcast(msg)
}

// SendChain sends the full chain to a specific peer.
func (s *P2PServer) SendChain(addr string, chain []blockchain.Block) {
	s.mu.RLock()
	p, ok := s.peers[addr]
	s.mu.RUnlock()

	if !ok {
		log.Printf("P2P: peer %s not found for SendChain", addr)
		return
	}

	msg := Message{Type: MsgChain, Chain: chain}
	if err := WriteMessage(p.conn, msg); err != nil {
		log.Printf("P2P: failed to send chain to %s: %v", addr, err)
	}
}

func (s *P2PServer) broadcast(msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for addr, p := range s.peers {
		if err := WriteMessage(p.conn, msg); err != nil {
			log.Printf("P2P: failed to send to %s: %v", addr, err)
		}
	}
}

// PeerCount returns the number of connected peers.
func (s *P2PServer) PeerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.peers)
}

// PeerAddresses returns the addresses of all connected peers.
func (s *P2PServer) PeerAddresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var addrs []string
	for addr := range s.peers {
		addrs = append(addrs, addr)
	}
	return addrs
}

// Stop shuts down the P2P server.
func (s *P2PServer) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.peers {
		p.conn.Close()
	}
	s.peers = make(map[string]*peer)
}
