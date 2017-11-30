package bitcoin

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/jmoiron/jsonq"
	stamp "github.com/daveqr/bitaudit/auditstamp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Message struct {
	method string
	txid   string
}

func GetBalance(config rpcclient.ConnConfig) (float64, error) {
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

func GetTransaction(config rpcclient.ConnConfig, txid string) interface{} {
	var result btcjson.GetTransactionResult
	if err := DoCommand(btcjson.NewGetTransactionCmd(txid, nil), &result, config); err != nil {
		panic(err)
	}

	return result
}

func SignMessage(request stamp.AuditStamp, config rpcclient.ConnConfig) (string, error) {
	//var result btcjson.SignMessageCmd
	//mki5RWvHxZHQQX6YY8Fm2ZnnYAMgzXxvvJ
	//'[{"txid":"4325a5db66cbc8e9ff6a585cd0e8a2288ea74f9b46d2972b93f63bbb7d09a23e","vout":0}]' '{"1AsJjnWg5QKBThM6mK9jZ8mmo6KUzDjRD":0.00186}'
	//var jsonStr = []byte(`{"method":"createrawtransaction", "params": [{"4325a5db66cbc8e9ff6a585cd0e8a2288ea74f9b46d2972b93f63bbb7d09a23e",0}] {"mki5RWvHxZHQQX6YY8Fm2ZnnYAMgzXxvvJ":1} }`)
	//
	//req, err := http.NewRequest("POST","http://" +  config.Host, bytes.NewBuffer(jsonStr))
	//req.SetBasicAuth(config.User, config.Pass)
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//body, _ := ioutil.ReadAll(resp.Body)
	//
	//data := map[string]interface{}{}
	//dec := json.NewDecoder(strings.NewReader(string(body)))
	//dec.Decode(&data)
	//jq := jsonq.NewQuery(data)
	//
	//fmt.Println(string(body))
	//return jq.String("result")
	var err = errors.New("asdf")
	return "", err
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

func ListTransactions(config rpcclient.ConnConfig) ([]map[string]interface{}, error) {
	log.Println("in list transactions")
	var jsonStr = []byte(`{"method":"listtransactions", "params": [""]}`)
	//var jsonStr = []byte(`{"method":"getpeerinfo", "params": [""]}`)

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
