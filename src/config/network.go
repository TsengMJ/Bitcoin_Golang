package config

import (
	"github.com/btcsuite/btcd/chaincfg"
)


func GetNetwork() *chaincfg.Params {
	config := config.GetString("network")

	if config == "mainnet" {
		return &chaincfg.MainNetParams
	} else if (config == "testnet") {
		return &chaincfg.TestNet3Params
	} else {
		return &chaincfg.RegressionNetParams
	}
}