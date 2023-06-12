package wallet

import (
	"encoding/hex"
	"log"

	"github.com/TsengMJ/Bitcoin_Golang/src/config"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)


type Wallet struct {
	privateKey 	*btcec.PrivateKey
	network		*chaincfg.Params
}

func Create() (Wallet, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Panic(err)
	}

	wallet := Wallet{
		privateKey: privateKey,
		network: config.GetNetwork(),
	}

	return wallet, nil
}


// Input Private Key is WIF format
func FromPrivateKey(wifPrivateKey string) (Wallet, error) {
	wif, err := btcutil.DecodeWIF(wifPrivateKey)
	if err != nil {
		log.Panic(err)
	}

	wallet := Wallet{
		privateKey: wif.PrivKey,
		network: config.GetNetwork(),
	}

	return wallet, nil
}

func (w *Wallet) GetPrivateKey() string {
	wif, err := btcutil.NewWIF(w.privateKey, w.network, true)
	if err != nil {
		log.Panic(err)
	}

	return wif.String()
}

func (w *Wallet) GetPublicKey() string {
	return hex.EncodeToString(w.privateKey.PubKey().SerializeCompressed())
}