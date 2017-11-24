package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"./commands"
)

func main() {
	StartServer()
}

type StampRequest struct {
	Message   string
	Signer   string
	Timestamp string
	Tx        string
}

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


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

func SignMessageCmd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	log.Println("r.PostForm", r.PostForm)

	decoder := schema.NewDecoder()
	stamp := new(StampRequest)
	err := decoder.Decode(stamp, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "sign message with:"+stamp.Message + stamp.Signer)
}

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", logger(HomeHandler)).Methods("GET")
	r.HandleFunc("/signmessage", logger(SignMessageCmd)).Methods("POST")
	r.HandleFunc("/signmessage", logger(HomeHandler)).Methods("Get")

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
