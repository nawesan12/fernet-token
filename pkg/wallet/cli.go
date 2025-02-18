package wallet

import (
	"bufio"
	"fmt"
	"os"
)

func CLI() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n🔹 Wallet CLI Menu")
		fmt.Println("1️⃣  Create a new wallet")
		fmt.Println("2️⃣  Show existing wallet")
		fmt.Println("3️⃣  Sign a transaction")
		fmt.Println("4️⃣  Exit")
		fmt.Print("Select an option: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			wallet, err := NewWallet()
			if err != nil {
				fmt.Println("❌ Failed to create wallet:", err)
				continue
			}
			fmt.Println("✅ Wallet Created!")
			fmt.Println("🔑 Public Key:", wallet.PublicKey)
			fmt.Println("🏦 Address:", wallet.Address)
		case "2":
			fmt.Println("📌 Feature not implemented yet!") // Can be linked to a persistent wallet store.
		case "3":
			fmt.Print("Enter transaction data to sign: ")
			scanner.Scan()
			data := scanner.Text()

			wallet, err := NewWallet() // This should ideally load an existing wallet.
			if err != nil {
				fmt.Println("❌ Error loading wallet:", err)
				continue
			}

			r, s, err := wallet.SignTransaction(data)
			if err != nil {
				fmt.Println("❌ Error signing transaction:", err)
				continue
			}

			fmt.Println("✅ Transaction Signed!")
			fmt.Println("🖋️  Signature (R):", r)
			fmt.Println("🖋️  Signature (S):", s)
		case "4":
			fmt.Println("👋 Exiting Wallet CLI.")
			return
		default:
			fmt.Println("⚠️ Invalid option, please try again.")
		}
	}
}
