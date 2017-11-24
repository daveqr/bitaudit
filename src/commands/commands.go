package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"strings"
	"log"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"encoding/hex"
)

const HashSize = 32
const MaxHashStringSize = HashSize * 2
type Hash [HashSize]byte

func  String(hash chainhash.Hash) string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

type Message struct {
	method string
	txid string
}

func GetBalance(config Config) (float64, error) {
	log.Println("in get balance")
	var jsonStr = []byte(`{"method":"getbalance"}`)

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(config.Username, config.Password)

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

func GetTransaction(config Config, txid string) (string, error) {

	hash := chainhash.DoubleHashH([]byte("def36856aa9fbebd028f12b7aed6d1f33b26758dae141adc79699fae052d6534"))
	log.Println(String(hash))
	log.Println(hex.EncodeToString(hash[:]))


	var jsonStr = []byte(`{"id":"gettransaction","jsonrpc":"2.0","method":"gettransaction", "params":["def36856aa9fbebd028f12b7aed6d1f33b26758dae141adc79699fae052d6534"]}`)

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(config.Username, config.Password)

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
	//jq := jsonq.NewQuery(data)

	fmt.Println(string(body))
	//return jq.String("result")
	return "ss", nil
}

func ListTransactions(config Config) ([]map[string]interface{}, error) {
	log.Println("in list transactions")
	var jsonStr = []byte(`{"method":"listtransactions", "params": [""]}`)
	//var jsonStr = []byte(`{"method":"getpeerinfo", "params": [""]}`)

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(config.Username, config.Password)

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
/*
func SignMessage(config Config, address string) (string, error) {
	var addr = string(ioutil.ReadFile("server.pem"))
	var jsonStr = []byte(`{"method":"signmessage", "params": ["addr", "hello world"]}`)

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(config.Username, config.Password)

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

	fmt.Println(string(body))
	return jq.String("result")
}
*/