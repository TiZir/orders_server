package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./interface/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", homePage).Methods("GET")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal(err)
	}
}
