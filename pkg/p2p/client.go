package p2p

import (
	"fmt"
	"net"
)

type P2PClient struct {
	ServerAddress string
}

func NewP2PClient(serverAddr string) *P2PClient {
	return &P2PClient{ServerAddress: serverAddr}
}

func (c *P2PClient) ConnectToPeer() {
	conn, err := net.Dial("tcp", c.ServerAddress)
	if err != nil {
		fmt.Println("❌ Failed to connect to peer:", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ Connected to peer:", c.ServerAddress)
}
