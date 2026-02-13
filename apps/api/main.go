package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nawesan12/fernet-token/packages/node"
)

func main() {
	httpPort := flag.String("http-port", "8080", "HTTP API port")
	p2pPort := flag.String("p2p-port", "6000", "P2P network port")
	dataDir := flag.String("data-dir", "", "Data directory (default: ~/.fernet-token)")
	miner := flag.String("miner", "", "Default miner address")
	peers := flag.String("peers", "", "Comma-separated list of seed peers (host:port)")
	flag.Parse()

	if *dataDir == "" {
		home, _ := os.UserHomeDir()
		*dataDir = home + "/.fernet-token"
	}
	os.MkdirAll(*dataDir, 0755)

	cfg := node.Config{
		DataDir: *dataDir,
		P2PPort: *p2pPort,
	}

	n, err := node.NewNode(cfg)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	// Start P2P
	n.StartP2P()

	// Connect to seed peers
	if *peers != "" {
		for _, addr := range strings.Split(*peers, ",") {
			addr = strings.TrimSpace(addr)
			if addr != "" {
				if err := n.P2P.ConnectToPeer(addr); err != nil {
					log.Printf("Failed to connect to peer %s: %v", addr, err)
				}
			}
		}
	}

	_ = miner // available for future auto-mining

	// Setup HTTP
	handler := NewAPIHandler(n)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + *httpPort,
		Handler: CORSMiddleware(JSONMiddleware(mux)),
	}

	go func() {
		log.Printf("HTTP API listening on port %s", *httpPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	n.Close()
	log.Println("Shutdown complete")
}
