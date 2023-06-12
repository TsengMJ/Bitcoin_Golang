package wallet

import (
	"log"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

type AddressType int
const (
	P2WPKH AddressType = iota
	P2PKH
	P2SH
	P2TR
)

func (w *Wallet) GetAddress(addressType AddressType) string {
	switch addressType {
	case P2WPKH:
		return w.GetAddressP2WPKH()
	case P2PKH:
		return w.GetAddressP2PKH()
	case P2SH:
		return w.GetAddressP2SH_P2WPKH()
	case P2TR:
		return w.GetAddressP2TR()
	default:
		return ""
	}
}

func (w *Wallet) GetAddressP2WPKH() string {
	p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(w.privateKey.PubKey().SerializeCompressed()), w.network)
	if err != nil {
		log.Fatal(err)
	}

	return p2wpkhAddress.EncodeAddress()
}

func (w *Wallet) GetAddressP2SH_P2WPKH() string {
	p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(w.privateKey.PubKey().SerializeCompressed()), w.network)
	if err != nil {
		log.Fatal(err)
	}

	// 將 P2WPKH 地址轉換為 P2SH 地址
	script, err := txscript.PayToAddrScript(p2wpkhAddress)
	if err != nil {
		log.Fatal(err)
	}

	p2shAddress, err := btcutil.NewAddressScriptHash(script, w.network)
	if err != nil {
		log.Fatal(err)
	}

	return p2shAddress.EncodeAddress()
}

func (w *Wallet) GetAddressP2TR() string {
	p2trAddress, _ := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(w.privateKey.PubKey())), w.network)

	return p2trAddress.EncodeAddress()
}

func (w *Wallet) GetAddressP2PKH() string {
	p2pkhAddress, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(w.privateKey.PubKey().SerializeCompressed()), w.network)
	if err != nil {
		log.Fatal(err)
	}

	return	p2pkhAddress.EncodeAddress()	
}