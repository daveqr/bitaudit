package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"net/url"
)

func logRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	logRequest(w, r)
	fmt.Fprintf(w, "Hello!")
}

func Xxx(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	logRequest(w, r)

	log.Println("r.PostForm", r.PostForm)
	log.Println("r.Form", r.Form)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("r.Body", string(body))

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("indices from body", values.Get("indices"))



	fmt.Fprintf(w, "xxxx!" + r.FormValue("data"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/xxx", Xxx)

	http.Handle("/", r)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
