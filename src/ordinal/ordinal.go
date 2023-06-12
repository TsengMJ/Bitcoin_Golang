package ordinal

import (
	"github.com/TsengMJ/Bitcoin_Golang/src/client"
	"github.com/TsengMJ/Bitcoin_Golang/src/model"
	"github.com/TsengMJ/Bitcoin_Golang/src/wallet"
	"github.com/btcsuite/btcd/wire"
)

type ContentType int

const (
	JSON ContentType = iota
	TEXT
)

type Inscriber struct {
	CommitFeeRate int
	RevealFeeRate int
	RevealOutValue int
}

const DEFAULT_INSCRIPTION_VALUE = 546

func NewInscriber(commitFeeRate int, revealFeeRate int, revealOutValue int) *Inscriber {
	return &Inscriber{
		CommitFeeRate: commitFeeRate,
		RevealFeeRate: revealFeeRate,
		RevealOutValue: revealOutValue,
	}
}

func (i *Inscriber) Inscribe (wallet *wallet.Wallet, addressType wallet.AddressType, content []byte, contentType ContentType) (*wire.MsgTx, *wire.MsgTx, error) {
	// Step 1: Select valid UTXOs
	apiClient := client.NewMempoolSpaceApiClent()
	utxos, err := apiClient.ListUnspent(wallet.GetAddress(addressType))
	if err != nil {
		return nil, nil, err
	}

	var selectedUtxos []model.UTXO
    for i := range utxos {
        if utxos[i].Value > DEFAULT_INSCRIPTION_VALUE {
            selectedUtxos = append(selectedUtxos, utxos[i])
        }
    }

	// Step 2: Create Commitment Transaction
	// Calculate the fee first
	commitTx := createCommitTx(selectedUtxos)

	// Step 3: Create Reveal Transaction
	revealTx := createRevealTx()


	return commitTx, revealTx, nil
}

func createCommitTx(utxos []model.UTXO ) (*wire.MsgTx) {
	commitTx := wire.NewMsgTx(wire.TxVersion)
	commitTx.TxIn = []*wire.TxIn{}

	return nil
}

func createRevealTx() (*wire.MsgTx) {
	return nil
}


func (i *Inscriber) GetCommitFee() uint64 {
	return i.CommitFeeRate
}

func (i *Inscriber) GetRevealFee() uint64 {
	return i.RevealFeeRate
}

func (i *Inscriber) GetRevealOutValue() uint64 {
	return i.RevealOutValue
}