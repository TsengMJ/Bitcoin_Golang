package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/TsengMJ/Bitcoin_Golang/src/config"
	"github.com/TsengMJ/Bitcoin_Golang/src/model"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// type MempoolSpaceApiClent ApiClient

type MempoolSpaceApiClent struct {
	baseUrl string
}

func NewMempoolSpaceApiClent() (*MempoolSpaceApiClent, error){
	newtwork := config.GetNetwork()

	switch newtwork.Name {
	case chaincfg.MainNetParams.Name:
		return &MempoolSpaceApiClent{
			baseUrl: "https://mempool.space/api",
		}, nil
	case chaincfg.TestNet3Params.Name:
		return &MempoolSpaceApiClent{
			baseUrl: "https://mempool.space/testnet/api",
		} , nil
	case chaincfg.RegressionNetParams.Name:
		return &MempoolSpaceApiClent{
			baseUrl: "https://mempool.space/signet/api",
		} , nil
	default:
		return &MempoolSpaceApiClent{
			baseUrl: "",
		} , errors.New("invalid network")
	}
}


func (c *MempoolSpaceApiClent) ListUnspent(address string) (*[]model.UTXO, error) {
	requestURL, err := url.JoinPath(c.baseUrl, "/address/", address, "/utxo")
	if err != nil {
		return nil, errors.New("failed to create request url: " + err.Error())
	}

	resBody, err := getRequest(requestURL)
	if err != nil {
		return nil, errors.New("failed to get request: " + err.Error())
	}

	var utxos []model.UTXO
	err = json.Unmarshal(resBody, &utxos)
	if err != nil {
		return nil, errors.New("failed to unmarshal response body: " + err.Error())
	}

	return &utxos, nil
}

func (c *MempoolSpaceApiClent) RecommendedFees () (*model.FeeRate, error) {
	requestURL, err := url.JoinPath(c.baseUrl, "/v1/fees/recommended")
	if err != nil {
		return nil, errors.New("failed to create request url: " + err.Error())
	}

	resBody, err := getRequest(requestURL)
	if err != nil {
		return nil, errors.New("failed to get request: " + err.Error())
	}

	var fees model.FeeRate
	err = json.Unmarshal(resBody, &fees)
	if err != nil {
		return nil, errors.New("failed to unmarshal response body: " + err.Error())
	}

	return &fees, nil
}

func (c *MempoolSpaceApiClent) BroadcastTx(tx *wire.MsgTx) (*chainhash.Hash, error) {
	requestURL, err := url.JoinPath(c.baseUrl, "/tx")
	if err != nil {
		return nil, errors.New("failed to create request url: " + err.Error())
	}

	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, errors.Errorf("failed to serialize tx: %s", err)
	}

	resBody, err := postRequest(requestURL, hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, errors.New("failed to post request: " + err.Error())
	}
	fmt.Println("resBody: ", string(resBody))

	txHash, err := chainhash.NewHashFromStr(string(resBody))
	fmt.Println("txHash: ", txHash)

	if err != nil {
		return nil, errors.Errorf("invalid tx hash: %s", err)
	}
	return txHash, nil
}


func getRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to get request: " + err.Error())
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	return resBody, nil
}

func postRequest(url string, hex string) ([]byte, error) {
	res, err := http.Post(url, "application/json", bytes.NewReader([]byte(hex)))
	if err != nil {
		return nil, errors.New("failed to post request: " + err.Error())
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	return resBody, nil
}