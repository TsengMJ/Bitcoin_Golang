package transaction

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type Transaction struct {
	tx *wire.MsgTx
	rpcclient *rpcclient.Client
	network *chaincfg.Params
}



func CreateTransaction(inputs []*wire.TxIn, outputs []*wire.TxOut) (*wire.MsgTx, error) {
	// 建立交易
	tx := wire.NewMsgTx(wire.TxVersion)

	tx.TxIn = inputs
	tx.TxOut = outputs

	return tx, nil
}

func (t *Transaction) Broadcast() error {
	return nil
}