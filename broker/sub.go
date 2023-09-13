package broker

import (
	"encoding/json"
	"log"
	"os"

	"github.com/TiZir/orders_server/db"
	"github.com/nats-io/nats.go"
)

func Subscribe() *nats.Conn {
	var err error

	conn, err := nats.Connect(os.Getenv("NT_URL"))
	if err != nil {
		log.Printf("error conn to nats: %v", err)
		return conn
	}
	_, err = conn.Subscribe("orders",
		func(msg *nats.Msg) {
			var order db.Order
			err := json.Unmarshal(msg.Data, &order)
			if err != nil {
				log.Printf("error unmarshal data to nats: %v\n", err)
				return
			}
			log.Println(order)

			_, err = db.AddOrder(order)
			if err != nil {
				log.Printf("error add order from nats: %v\n", err)
				return
			}

			log.Printf("received order: %v\n", order.OrderUID)
		},
	)
	if err != nil {
		log.Printf("error subscribe nats: %v", err)
	}
	return conn
}
