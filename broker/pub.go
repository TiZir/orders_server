package broker

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/TiZir/orders_server/db"
	"github.com/nats-io/nats.go"
)

func Publish() *nats.Conn {
	item1 := db.Item{
		ChrtID: 1, TrackNumber: "WBILMTESTTRACKRR", Price: 10, Rid: "RID 1", Name: "T-Shirt",
		Sale: 9, Size: "S", TotalPrice: 13, NmID: 1, Brand: "Adidas", Status: 1,
	}
	item2 := db.Item{
		ChrtID: 2, TrackNumber: "WBILMTESTTRACKRR", Price: 11, Rid: "RID 2", Name: "T-Shirt",
		Sale: 10, Size: "L", TotalPrice: 14, NmID: 2, Brand: "Nike", Status: 1,
	}
	item3 := db.Item{
		ChrtID: 3, TrackNumber: "WBILMTESTTRACKRR", Price: 12, Rid: "RID 3", Name: "T-Shirt",
		Sale: 11, Size: "M", TotalPrice: 15, NmID: 3, Brand: "Polo", Status: 1,
	}

	delivery := db.Delivery{
		Name: "gruzovichki", Phone: "+70000000000", Zip: "2639809", City: "Pupok Sever",
		Address: "Nijnee Puzo", Region: "Vishe Taza", Email: "bobus@mail.ru",
	}

	payment := db.Payment{
		Transaction: "testtrantest", RequestID: "", Currency: "USD", Provider: "WBPay",
		Amount: 1703, PaymentDt: 2, Bank: "SBER", DeliveryCost: 150, GoodsTotal: 3, CustomFee: 4,
	}

	order := db.Order{
		OrderUID: "testtrantest", TrackNumber: "WBILMTESTTRACKRR", Entry: "ent",
		Delivery: delivery, Payment: payment, Items: []db.Item{item1, item2, item3},
		Locale: "RU", InternalSignature: "", CustomerID: "test",
		DeliveryService: "vanomasENT", Shardkey: "testkey", SmID: 5,
		DateCreated: "2021-11-26T06:22:19Z", OofShard: "1",
	}

	orderData, err := json.Marshal(order)
	if err != nil {
		log.Printf("error marshal data to nats: %v\n", err)
	}

	conn, err := nats.Connect(os.Getenv("NT_URL"))
	if err != nil {
		log.Printf("error conn to nats: %v\n", err)
		return conn
	}
	time.Sleep(5 * time.Second)
	err = conn.Publish("orders", orderData)
	if err != nil {
		log.Printf("error nats publish: %v\n", err)
		return conn
	}
	return conn
}
