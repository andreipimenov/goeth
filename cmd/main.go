package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/andreipimenov/goeth"
)

type Account struct {
	Address    string
	PrivateKey string
}

func main() {

	// Two accounts created previously.
	account1 := Account{
		Address:    "0x6A7999bF334185DC0E81B37Ab3a53F3CcEc85ee6",
		PrivateKey: "482066177d9ac64c66f6d54e4800896513017567d45707cacb372ef137a07729",
	}
	account2 := Account{
		Address:    "0x01CF9aD892602D94f21772cC6b3C7fd229b60951",
		PrivateKey: "5b3b197a6bcc3993d35292f57f60a05e4b522985ac57695a08f382f795fe6748",
	}

	// Connecting to Infura test network.
	client, err := goeth.ConnectToInfura(goeth.InfuraRopstenNet)
	if err != nil {
		log.Fatalf("Error connecting to infura network: %v\n", err)
	}

	ctx := context.Background()

	// Get current balance of the accounts.
	balance1, err := goeth.Balance(ctx, client, account1.Address)
	if err != nil {
		log.Fatalf("Error getting balance of %s: %v\n", account1.Address, err)
	}
	fmt.Printf("Balance of %s: %v ETH\n", account1.Address, balance1)

	balance2, err := goeth.Balance(ctx, client, account2.Address)
	if err != nil {
		log.Fatalf("Error getting balance of %s: %v\n", account2.Address, err)
	}
	fmt.Printf("Balance of %s: %v ETH\n", account2.Address, balance2)

	// Create transaction for sending 0.1 ETH from account 1 to account 2.
	tx, err := goeth.NewTx(ctx, client, account1.Address, account2.Address, big.NewInt(100000000000000000))
	if err != nil {
		log.Fatalf("Error creating transaction: %v\n", err)
	}

	// Sign transaction.
	signedTx, err := goeth.SignTx(ctx, client, tx, account1.PrivateKey)
	if err != nil {
		log.Fatalf("Error signing transaction: %v\n", err)
	}

	// Send transation.
	err = goeth.SendTx(ctx, client, signedTx)
	if err != nil {
		log.Fatalf("Error sending transaction: %v\n", err)
	}

	fmt.Printf("Transaction %s sent\n", signedTx.Hash().Hex())

	// Start wainting for transation being mined.
	for {
		fmt.Println("Waiting for 5 seconds until block being mined...")
		<-time.After(5 * time.Second)
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil && receipt.Status == 1 {
			fmt.Printf("Receipt status for %s: Success\n", receipt.TxHash.Hex())
			break
		}
	}
}
