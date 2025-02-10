package p2p

import (
	"fmt"
	"net"
)

type P2PServer struct {
	Peers []string
}

func NewP2PServer() *P2PServer {
	return &P2PServer{Peers: []string{}}
}

func (s *P2PServer) Start(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("🔗 P2P Server listening on port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *P2PServer) handleConnection(conn net.Conn) {
	fmt.Println("New peer connected:", conn.RemoteAddr().String())
	conn.Close()
}
