package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/TiZir/orders_server/broker"
	"github.com/TiZir/orders_server/cash"
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
	ch := cash.GetCash()
	var order db.Order
	var err error

	data, ok := ch.Load(id)
	if ok {
		order, ok = data.(db.Order)
		if !ok {
			log.Printf("error type assertion from cash: %v\n", err)
			return
		}
		log.Println("get from cash")
	} else {
		order, err = db.SelectOrderByID(id)
		if err != nil {
			log.Printf("incorrect id for order: %v\n", err)
			return
		}
		log.Println("get from bd")
	}

	orderData, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(orderData)
}

// func addOrder(w http.ResponseWriter, r *http.Request) {
// 	var data db.Order
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	ok, err := db.AddOrder(data)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	fmt.Fprint(w, ok)
// }

func main() {
	router := mux.NewRouter()
	conn := broker.Subscribe()
	defer conn.Close()
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/orders/{id}", orderPage).Methods("GET")
	// router.HandleFunc("/orders/add", addOrder).Methods("POST")
	if err := http.ListenAndServe(":9090", router); err != nil {
		log.Fatal(err)
	}
}
