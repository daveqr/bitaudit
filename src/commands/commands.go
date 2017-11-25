package commands

import (
	"bytes"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	method string
	txid   string
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

func GetTransaction(config Config, txid string) interface{} {
	var result btcjson.GetTransactionResult
	if err := DoCommand(btcjson.NewGetTransactionCmd(txid, nil), &result, config); err != nil {
		panic(err)
	}

	return result
}

func DoCommand(cmd interface{}, result interface{}, config Config) (error) {
	id := 1
	marshalledBytes, err := btcjson.MarshalCmd(id, cmd)
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", config.Url, bytes.NewBuffer(marshalledBytes))
	req.SetBasicAuth(config.Username, config.Password)

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
