package blockchain

import "fmt"

// SimpleConsensus verifica si la transacción es válida
func SimpleConsensus(tx *Transaction) bool {
	if tx.Amount <= 0 {
		fmt.Println("Transacción inválida: monto no válido")
		return false
	}
	return true
}
