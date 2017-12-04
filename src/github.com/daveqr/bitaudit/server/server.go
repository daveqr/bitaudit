package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

func logger(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("start request")
		log.Println(r.Form)
		log.Println("method", r.Method)
		log.Println("path", r.URL.Path)
		log.Println("scheme", r.URL.Scheme)
		log.Println(r.Form["url_long"])
		for k, v := range r.Form {
			log.Println("key:", k)
			log.Println("val:", strings.Join(v, ""))
		}

		fn(w, r)
		log.Println("finish request")
	}
}

func StartServer() {
	log.Println("Starting server")

	r := mux.NewRouter()

	r.HandleFunc("/", logger(HomeHandler)).Methods("GET")
	r.HandleFunc("/signmessage", logger(SignMessageHandler)).Methods("POST")
	r.HandleFunc("/signmessage", logger(HomeHandler)).Methods("Get")
	r.HandleFunc("/verifymessage/{signer}/{txId}", logger(VerifyMessageHandler)).Methods("Get")
	r.HandleFunc("/checkbalance", logger(BalanceHandler)).Methods("Get")
	r.HandleFunc("/checkbitcoinbalance", logger(BitcoinBalanceHandler)).Methods("Get")
	r.HandleFunc("/listbitcointransactions", logger(BitcoinListTransactionsHandler)).Methods("Get")
	r.HandleFunc("/printlocalchain", logger(PrintLocalChainHandler)).Methods("Get")

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
