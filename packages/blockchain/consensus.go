package blockchain

import "fmt"

// SimpleConsensus valida una transacción antes de agregarla a la cadena
func SimpleConsensus(tx *Transaction) bool {
	if tx == nil {
		fmt.Println("❌ Transacción inválida: es nula")
		return false
	}
	if !tx.IsValid() {
		fmt.Println("❌ Transacción inválida: datos incorrectos")
		return false
	}
	return true
}
