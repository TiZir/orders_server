package broker

import (
	"encoding/json"
	"log"

	"github.com/TiZir/orders_server/db"
	"github.com/nats-io/nats.go"
)

func Subscribe() (*nats.Conn, error) {
	var err error

	conn, err := nats.Connect("nats://nats-streaming:4222")
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Subscribe("orders",
		func(msg *nats.Msg) {
			var order db.Order
			err := json.Unmarshal(msg.Data, &order)
			if err != nil {
				log.Printf("Error unmarshaling order: %v\n", err)
				return
			}
			log.Println(order)

			_, err = db.AddOrder(order)
			if err != nil {
				log.Printf("Error unmarshaling order: %v\n", err)
				return
			}

			log.Printf("Received order: %v\n", order.OrderUID)
		},
	)
	if err != nil {
		log.Printf("error subscribe nats: %v", err)
		return conn, err
	}
	return conn, nil
}

// func messageHandler(data []byte) bool {
// 	recievedOrder := db.Order{}
// 	err := json.Unmarshal(data, &recievedOrder)
// 	if err != nil {
// 		log.Printf("messageHandler() error, %v\n", err)
// 		return true
// 	}

// 	_, err = s.dbObject.AddOrder(recievedOrder)
// 	if err != nil {
// 		log.Printf("%s: unable to add order: %v\n", s.name, err)
// 		return false
// 	}
// 	return true
// }
