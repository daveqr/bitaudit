package server

import (
	"fmt"
	"github.com/daveqr/bitaudit/blockchain"
	stamp "github.com/daveqr/bitaudit/auditstamp"
	"github.com/gorilla/schema"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"encoding/json"
)

var commands blockchain.Commands

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	address := r.URL.Query().Get("address")

	balance := commands.GetBalance(address)

	log.Println("Balance: " + fmt.Sprintf("%f", balance))

	t, _ := template.ParseFiles("templates/balance.html")
	t.Execute(w, balance)
}

func VerifyMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	address := r.URL.Query().Get("address")

	balance := commands.GetBalance(address)

	log.Println("Balance: " + fmt.Sprintf("%f", balance))

	t, _ := template.ParseFiles("templates/balance.html")
	t.Execute(w, balance)
}

func PrintLocalChainHandler(w http.ResponseWriter, r *http.Request) {
	commands.PrintChain()
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func SignMessageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	log.Println("r.PostForm", r.PostForm)

	auditStamp, err := createStamp(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx := commands.WriteMessageToBlockchain(auditStamp.Hash().String())
	auditStamp.Txid = fmt.Sprintf("%x", tx.ID)
	stamp.SaveToDb(auditStamp)

	var tmp struct {
		Txid string `json:"Txid"`
	}
	tmp.Txid = auditStamp.Txid
	a, err := json.Marshal(tmp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(a)
}

func createStamp(r *http.Request) (*stamp.AuditStamp, error) {
	auditStamp := new(stamp.AuditStamp)

	decoder := schema.NewDecoder()
	err := decoder.Decode(auditStamp, r.PostForm)
	if err != nil {
		return nil, err
	}

	auditStamp.Timestamp = time.Now()
	auditStamp.Status = stamp.StatusCreated

	return auditStamp, nil
}

func VerifyMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txId := vars["txId"]

	as := stamp.GetFromDb(txId)

	w.Write([]byte("Status: " + as.StatusString()))
}
