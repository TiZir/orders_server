package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d Delivery) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Delivery) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &d)
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

func (p Payment) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Payment) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &p)
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (i Item) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *Item) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &i)
}

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             Items     `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}
type Items []Item

func (is *Items) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("error casting value to bytes")
	}

	var items []Item
	err := json.Unmarshal(bytes, &items)
	if err != nil {
		return err
	}

	*is = items
	return nil
}

func (is Items) Value() (driver.Value, error) {
	return json.Marshal(is)
}
