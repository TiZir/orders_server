package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/TiZir/orders_server/broker"
	"github.com/TiZir/orders_server/db"
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

func orderPage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	order, err := db.SelectOrderByID(id)
	if err != nil {
		log.Printf("Incorrect id for order: %v", err)
		return
	}
	tmpl := template.Must(template.ParseFiles("./interface/index.html"))
	// orderData, err := json.Marshal(order)
	// if err != nil {
	// 	log.Println(err)
	// }
	err = tmpl.Execute(w, map[string]interface{}{
		"orderData": order,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func addOrder(w http.ResponseWriter, r *http.Request) {
	var data db.Order
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		return
	}
	ok, err := db.AddOrder(data)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprint(w, ok)
}

func main() {
	router := mux.NewRouter()
	conn, err := broker.Subscribe()
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/orders/{id}", orderPage).Methods("GET")
	router.HandleFunc("/orders/add", addOrder).Methods("POST")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal(err)
	}
}
