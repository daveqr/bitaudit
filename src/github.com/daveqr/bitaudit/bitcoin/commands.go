package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	stamp "github.com/daveqr/bitaudit/auditstamp"
	"github.com/daveqr/bitaudit/writebtc"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Message struct {
	method string
	txid   string
}

func GetWalletBalance() (float64, error) {
	log.Println("in get balance")
	var jsonStr = []byte(`{"method":"getbalance"}`)

	req, err := http.NewRequest("POST", "http://"+config.Host, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(config.User, config.Pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq.Float("result")
}

func GetTransaction(txid string) interface{} {
	var result btcjson.GetTransactionResult
	if err := DoCommand(btcjson.NewGetTransactionCmd(txid, nil), &result, config); err != nil {
		panic(err)
	}

	return result
}

func SignMessage(as *stamp.AuditStamp) (*chainhash.Hash) {

	payToAddr := "mm4dSgXb1WwSamdQFsszkmvXpXCNyPQ5h1"

	address, err := btcutil.DecodeAddress(payToAddr, writebtc.ActiveNet.Params)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	asHash := as.Hash()
	opReturn, err := txscript.NullDataScript(asHash.CloneBytes())
	payToScript, err := txscript.PayToAddrScript(address)

	ntfnHandlers := rpcclient.NotificationHandlers{
		// TODO check for acceptance
	}

	msgtx := wire.NewMsgTx(1)

	msgtx.AddTxOut(wire.NewTxOut(0, opReturn))
	msgtx.AddTxOut(wire.NewTxOut(int64(1), payToScript))

	clientx, err := rpcclient.New(&config, &ntfnHandlers)
	h, err := clientx.SendRawTransaction(msgtx, true)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return h
}

func Encode(hash []byte) ([]byte, error) {
	length := len(hash)
	if length == 0 || length > 30 {
		return nil, errors.New("Encode can only handle 1 to 30 bytes")
	}
	var b []byte = make([]byte, 0, 33)
	b = append(b, byte(2), byte(length))

	b = append(b, hash...)

	if length < 30 {
		data := make([]byte, 30-length)
		b = append(b, data...)
	}

	b = append(b, byte(0))

	for i := 0; i < 256; i++ {
		b[len(b)-1] = byte(i)
		adr2 := hex.EncodeToString(b)
		_, e := btcutil.DecodeAddress(adr2, writebtc.ActiveNet.Params)
		if e == nil {
			return b, nil
		}
	}

	log.Print("Failure")
	return b, errors.New("Couldn't fix the address")
}

func Decode(addr []byte) []byte {
	length := int(addr[1])
	data := addr[2 : length+2]
	return data
}

func DoCommand(cmd interface{}, result interface{}, config rpcclient.ConnConfig) error {
	id := 1
	marshalledBytes, err := btcjson.MarshalCmd(id, cmd)
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", "http://"+config.Host, bytes.NewBuffer(marshalledBytes))
	req.SetBasicAuth(config.User, config.Pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	buf := new(bytes.Buffer)
	buf.ReadFrom(strings.NewReader(string(body)))

	var response btcjson.Response
	if err := json.Unmarshal([]byte(buf.String()), &response); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(response.Result, &result); err != nil {
		panic(err)
	}

	return nil
}

func ListTransactions() ([]map[string]interface{}, error) {
	log.Println("in list transactions")
	var jsonStr = []byte(`{"method":"listtransactions", "params": [""]}`)

	req, err := http.NewRequest("POST", "http://"+config.Host, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(config.User, config.Pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	return jq.ArrayOfObjects("result")
}