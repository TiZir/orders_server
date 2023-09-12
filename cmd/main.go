package main

import (
	"log"

	"github.com/TiZir/orders_server/broker"
)

func main() {
	conn, err := broker.Publish()
	if err != nil {
		log.Println(err)
	}
	err = conn.Publish("orders", []byte("trash"))
	if err != nil {
		log.Println("ERROR: conn.Publish:", err)
	}
	defer conn.Close()
}
