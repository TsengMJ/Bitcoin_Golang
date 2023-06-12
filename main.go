package main

import (
	"fmt"

	"github.com/TsengMJ/Bitcoin_Golang/src/client"
	"github.com/TsengMJ/Bitcoin_Golang/src/config"
)


func main() {
	config.Init("development")

	/* Wallet Example */
	// wallet, _ := wallet.FromPrivateKey("cRBVmt1k3A7cggpxgwexfSYawBsGFiCRzQFpeoV3i5URxrTRmMsi")

	// println("Private Key: ", wallet.GetPrivateKey())
	// println("Public Key: ", wallet.GetPublicKey())
	// println("Address P2WPKH: ", wallet.GetAddressP2WPKH())
	// println("Address P2SH: ", wallet.GetAddressP2SH_P2WPKH())
	// println("Address P2TR: ", wallet.GetAddressP2TR())
	// println("Address P2PKH: ", wallet.GetAddressP2PKH())


	/* Transaction Example */
	client := client.NewMempoolSpaceApiClent()
	utxos, err := client.ListUnspent("tb1q4kgratttzjvkxfmgd95z54qcq7y6hekdm3w56u")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", utxos)
}