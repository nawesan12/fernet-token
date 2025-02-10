package p2p

import (
	"fmt"
	"net"
)

type P2PClient struct{}

func NewP2PClient() *P2PClient {
	return &P2PClient{}
}

func (c *P2PClient) ConnectToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Failed to connect to peer:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to peer:", address)
}
