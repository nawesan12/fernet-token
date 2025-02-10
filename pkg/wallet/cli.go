package wallet

import "fmt"

func CLI() {
	wallet := NewWallet()
	if wallet == nil {
		fmt.Println("Failed to create wallet")
		return
	}
	fmt.Println("🪙 New Wallet Created!")
	fmt.Println("Public Key:", wallet.PublicKey)
}
