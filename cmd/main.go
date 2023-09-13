package main

import (
	"log"

	"github.com/TiZir/orders_server/broker"
)

func main() {
	conn := broker.Publish()

	//trash data test
	err := conn.Publish("orders", []byte("trash"))
	if err != nil {
		log.Printf("error nats publish: %v\n", err)
	}

	defer conn.Close()
}
