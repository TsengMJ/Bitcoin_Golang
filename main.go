package main

import (
	"fmt"
	"os"

	"github.com/TsengMJ/Bitcoin_Golang/src/config"
	"github.com/TsengMJ/Bitcoin_Golang/src/wallet"
)


func main() {
	config.Init("development")

	/* Wallet Example */
	testWallet, _ := wallet.FromPrivateKey("cPcvGhULd5FcdS4EuMngs83HjiAdqnTQffcDfaufp9JvKDn9UgmM")

	receiveAddr, err := testWallet.GetAddressP2WPKH()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	txHash, err := testWallet.Send(wallet.P2WPKH, receiveAddr, 1000, 3)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("TxHash:", txHash)
}