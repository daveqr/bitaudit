package commands

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"strings"
	"encoding/json"
	"github.com/jmoiron/jsonq"
	"fmt"
)

func GetBalance(config Config) (float64, error) {
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
	var jsonStr = []byte(`{"method":"gettransaction", "params": ["txid"]}`)

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

func ListTransactions(config Config) (float64, error) {
	var jsonStr = []byte(`{"method":"listtransactions", "params": [""]}`)

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
