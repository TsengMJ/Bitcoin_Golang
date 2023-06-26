package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/TsengMJ/Bitcoin_Golang/src/client"
	"github.com/TsengMJ/Bitcoin_Golang/src/config"
	"github.com/TsengMJ/Bitcoin_Golang/src/model"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)


type Wallet struct {
	privateKey 	*btcec.PrivateKey
	network		*chaincfg.Params
}

func Create() (Wallet, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return Wallet{}, errors.New("failed to create private key with btece: " + err.Error())
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
		return Wallet{}, errors.New("failed to decode WIF private key: " + err.Error())
	}

	wallet := Wallet{
		privateKey: wif.PrivKey,
		network: config.GetNetwork(),
	}

	return wallet, nil
}

func (w *Wallet) GetPrivateKey() (string, error) {
	wif, err := btcutil.NewWIF(w.privateKey, w.network, true)
	if err != nil {
		return "", errors.New("failed to convert private key to WIF: " + err.Error())
	}

	return wif.String(), nil
}

func (w *Wallet) GetPublicKey() string {
	return hex.EncodeToString(w.privateKey.PubKey().SerializeCompressed())
}

func (w *Wallet) Send(sendAddressType AddressType, recipient string, amount int64, feeRate int64) (string, error) {
	// Get Client
	client, err := client.NewMempoolSpaceApiClent()
	if err != nil {
		return "", errors.New("failed to create mempool space api client: " + err.Error())
	}

	// Get Send Address
	sendAddress, err := w.GetAddress(sendAddressType)
	if err != nil {
		return "", errors.New("failed to get send address: " + err.Error())
	}

	// Get UTXOs
	utxos, err := client.ListUnspent(sendAddress);
	if err != nil {
		return "", errors.New("failed to get utxos: " + err.Error())
	}

	// Filterout Small UTXOs
	utxos, err = filterSmallUTXOs(utxos)
	if err != nil {
		return "", errors.New("failed to filter small utxos: " + err.Error())
	}

	// Create Empty Transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Prepare TxInputs
	for _, utxo := range *utxos {
		hash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return "", errors.New("failed to convert txid to hash: " + err.Error())
		}
		outPoint := wire.NewOutPoint(hash, utxo.Vout)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	// Prepare recepient output
	rcvAddress, err := btcutil.DecodeAddress(recipient, config.GetNetwork())
	if err != nil {
		return "", errors.New("failed to decode recipient address: " + err.Error())
	}

    rcvScript, err := txscript.PayToAddrScript(rcvAddress)
	if err != nil {
		return "", errors.New("failed to create recipient script: " + err.Error())
	}

	rcvOut := wire.NewTxOut(amount, rcvScript)
	tx.AddTxOut(rcvOut)

	// Prepare change output
	changeAddress, err := btcutil.DecodeAddress(sendAddress, config.GetNetwork())
	if err != nil {
		return "", errors.New("failed to decode change address: " + err.Error())
	}

	changeScript, err := txscript.PayToAddrScript(changeAddress)
	if err != nil {
		return "", errors.New("failed to create change script: " + err.Error())
	}

	changeOutput := wire.NewTxOut(0, changeScript)
	tx.AddTxOut(changeOutput)

	// Set Tx Fee
	var totalUTXOValue int64
	for _, utxo := range *utxos {
		totalUTXOValue += utxo.Value
	}

	txVirtualSize := mempool.GetTxVirtualSize(btcutil.NewTx(tx))
	txFee := btcutil.Amount(txVirtualSize) * btcutil.Amount(feeRate)
	changeAmount := btcutil.Amount(totalUTXOValue) - btcutil.Amount(amount) - txFee
	tx.TxOut[len(tx.TxOut) - 1].Value = int64(changeAmount)

	fmt.Println("TxInput Length: ", len(tx.TxIn))


	// Sign TxInputs
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	// prevOutputs := []wire.TxOut{}
	for i, txIn := range tx.TxIn {
		prevOutput := wire.NewTxOut((*utxos)[i].Value, changeScript)
		fmt.Printf( "prevOutput: %v\n", prevOutput )
		prevOutFetcher.AddPrevOut(txIn.PreviousOutPoint, prevOutput)
	}

	for i, txIn := range tx.TxIn {
		// prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
		// prevOutput := wire.NewTxOut((*utxos)[i].Value, changeScript)
		// fmt.Printf( "prevOutput: %v\n", prevOutput )
		// prevOutFetcher.AddPrevOut(txIn.PreviousOutPoint, prevOutput)

		var witness wire.TxWitness
		prevOutput := prevOutFetcher.FetchPrevOutput(txIn.PreviousOutPoint)

		switch changeAddress.(type){
		case *btcutil.AddressTaproot:
			witness, err = txscript.TaprootWitnessSignature(
				tx,
				txscript.NewTxSigHashes(tx, prevOutFetcher),
				i,
				prevOutput.Value,
				prevOutput.PkScript,
				txscript.SigHashAll,
				w.privateKey,
			)
			if err != nil {
				return "", errors.New("failed to sign taproot tx: " + err.Error())
			}
		default:
			witness, err = txscript.WitnessSignature(
				tx, 
				txscript.NewTxSigHashes(tx, prevOutFetcher), 
				i, 
				prevOutput.Value, 
				prevOutput.PkScript, 
				txscript.SigHashAll, 
				w.privateKey, 
				true,
			)
			if err != nil {
				return "", errors.New("failed to sign witness tx: " + err.Error())
			}
		}

		tx.TxIn[i].Witness = witness
	}


	// Broadcast Tx
	txHash, err := client.BroadcastTx(tx)
	if err != nil {
		return "", errors.New("failed to broadcast tx: " + err.Error())
	}

	return txHash.String(), nil	
}

func filterSmallUTXOs(utxos *[]model.UTXO) (*[]model.UTXO,  error) {
	var filteredUTXOs []model.UTXO

	for _, utxo := range *utxos {
		if utxo.Value > 546 {
			filteredUTXOs = append(filteredUTXOs, utxo)
		}
	}

	return &filteredUTXOs, nil
}