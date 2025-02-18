package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

// Estructura de un bloque
type Block struct {
	Index        int
	Timestamp    string
	PrevHash     string
	Hash         string
	Transactions []string
}

// Genera un hash único basado en los datos del bloque
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.Index, block.Timestamp, block.PrevHash, block.Transactions)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Crea un nuevo bloque
func createBlock(prevBlock Block, transactions []string) Block {
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		PrevHash:     prevBlock.Hash,
		Transactions: transactions,
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func main() {
	// Bloque génesis
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		PrevHash:     "0",
		Transactions: []string{"Genesis Block"},
	}
	genesisBlock.Hash = calculateHash(genesisBlock)

	// Crear un segundo bloque
	newBlock := createBlock(genesisBlock, []string{"Pago de 100 Fernet", "Pago de 50 Fernet"})
	fmt.Println("Bloque creado:", newBlock)

	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()

	fmt.Println("Nodo escuchando en el puerto 8080...")

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Nueva conexión establecida:", conn.RemoteAddr())
}
