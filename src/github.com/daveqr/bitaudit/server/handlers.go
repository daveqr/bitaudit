package server

import (
	"fmt"
	stamp "github.com/daveqr/bitaudit/auditstamp"
	"github.com/daveqr/bitaudit/bitcoin"
	"github.com/daveqr/bitaudit/blockchain"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"log"
	"net/http"
	"time"
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

func BitcoinBalanceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	balance, _ := bitcoin.GetWalletBalance()

	log.Println("Bitcoin Wallet Balance: " + fmt.Sprintf("%f", balance))

	t, _ := template.ParseFiles("templates/balance.html")
	t.Execute(w, balance)
}

func BitcoinListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	balance, _ := bitcoin.ListTransactions()

	log.Println("Bitcoin Wallet Balance: " + fmt.Sprintf("%f", balance))

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

	as, err := createStamp(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//bitcoin.SignMessage(as)
	//if h != nil {
	//	log.Println("h is not null")
	//}

	tx := commands.WriteMessageToBlockchain(as.Hash().String())
	as.Txid = fmt.Sprintf("%x", tx.ID)

	stamp.SaveToDb(as)
	retJson, err := as.ReturnJson()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(retJson)
}

func createStamp(r *http.Request) (*stamp.AuditStamp, error) {
	auditStamp := new(stamp.AuditStamp)

	decoder := schema.NewDecoder()
	err := decoder.Decode(auditStamp, r.PostForm)
	if err != nil {
		return nil, err
	}

	auditStamp.Message = r.FormValue("message")
	auditStamp.Signer = r.FormValue("signer")
	auditStamp.Timestamp = time.Now()
	auditStamp.Status = stamp.StatusCreated

	return auditStamp, nil
}

func VerifyMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tempAs := stamp.AuditStamp{Signer: vars["signer"], Txid: vars["txId"]}

	as := stamp.GetFromDb(tempAs.Key())

	w.Write([]byte("Status: " + as.StatusString()))
}
