package wallet

import (
	"errors"

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

func (w *Wallet) GetAddress(addressType AddressType) (string, error) {
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
		return "", errors.New("unknown address type")
	}
}

func (w *Wallet) GetAddressP2WPKH() (string, error) {
	p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(w.privateKey.PubKey().SerializeCompressed()), w.network)
	if err != nil {
		return "", errors.New("failed to create P2WPKH address: " + err.Error())
	}

	return p2wpkhAddress.EncodeAddress(), nil
}

func (w *Wallet) GetAddressP2SH_P2WPKH() (string, error)  {
	p2wpkhAddress, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(w.privateKey.PubKey().SerializeCompressed()), w.network)
	if err != nil {
		return "", errors.New("failed to create P2WPKH address: " + err.Error())
	}

	// 將 P2WPKH 地址轉換為 P2SH 地址
	script, err := txscript.PayToAddrScript(p2wpkhAddress)
	if err != nil {
		return "", errors.New("failed to convert P2WPKH address to P2SH address: " + err.Error())
	}

	p2shAddress, err := btcutil.NewAddressScriptHash(script, w.network)
	if err != nil {
		return "", errors.New("failed to create P2SH address: " + err.Error())
	}

	return p2shAddress.EncodeAddress(), nil
}

func (w *Wallet) GetAddressP2TR() (string, error)  {
	p2trAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(w.privateKey.PubKey())), w.network)
	if err != nil {
		return "", errors.New("failed to create P2TR address: " + err.Error())
	}

	return p2trAddress.EncodeAddress(), nil
}

func (w *Wallet) GetAddressP2PKH() (string, error)  {
	p2pkhAddress, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(w.privateKey.PubKey().SerializeCompressed()), w.network)
	if err != nil {
		return "", errors.New("failed to create P2PKH address: " + err.Error())
	}

	return	p2pkhAddress.EncodeAddress(), nil
}