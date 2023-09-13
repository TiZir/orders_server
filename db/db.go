package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/TiZir/orders_server/cash"
)

func SelectOrderByID(id string) (Order, error) {
	var o Order
	var db *sql.DB
	var err error
	ch := cash.GetCash()

	db, err = GetDB()
	if err != nil {
		log.Printf("error get db: %v\n", err)
		return o, err
	}

	//Order
	err = db.QueryRow(`SELECT order_uid, track_number, entry, locale, COALESCE(internal_signature, ''), customer_id, 
	delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1;`, id).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
		&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard,
	)
	if err != nil {
		log.Println(id, o)
		log.Printf("error select order: %v\n", err)
		return o, err
	}

	//Payment
	err = db.QueryRow(`SELECT transaction, COALESCE(request_id, ''), currency, provider, amount, payment_dt, bank,
	delivery_cost, goods_total, custom_fee FROM payment WHERE transaction = $1;`, o.OrderUID).Scan(
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider,
		&o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)
	if err != nil {
		log.Println(id, o)
		log.Printf("error select payment: %v\n", err)
		return o, err
	}

	//Delivery
	err = db.QueryRow(`SELECT name, phone, zip, city, address, region,
	email FROM delivery WHERE order_uid = $1;`, o.OrderUID).Scan(
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip,
		&o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region,
		&o.Delivery.Email,
	)
	if err != nil {
		log.Println(id, o)
		log.Printf("error select delivery: %v\n", err)
		return o, err
	}

	//Items
	rows, err := db.Query(`SELECT * FROM items WHERE track_number = $1;`, o.TrackNumber)
	if err != nil {
		log.Println(id, o)
		log.Printf("error select items: %v\n", err)
		return o, err
	}
	defer rows.Close()

	for rows.Next() {
		var i Item
		err = rows.Scan(
			&i.ChrtID, &i.TrackNumber, &i.Price, &i.Rid, &i.Name,
			&i.Sale, &i.Size, &i.TotalPrice, &i.NmID,
			&i.Brand, &i.Status,
		)
		if err != nil {
			log.Println(id, o)
			log.Printf("error scan items: %v\n", err)
			return o, err
		}
		o.Items = append(o.Items, i)
	}
	ch.Store(o.OrderUID, o)
	return o, nil
}

func AddOrder(o Order) (string, error) {
	var db *sql.DB
	var err error
	//ch := cash.GetCash()
	db, err = GetDB()
	if err != nil {
		log.Printf("error open db: %v\n", err)
		return "", err
	}
	tx, err := db.Begin()
	if err != nil {
		log.Printf("error transaction start: %v\n", err)
		return "", err
	}
	//Order
	var internalSignature interface{}
	if o.InternalSignature == "" {
		internalSignature = nil
	} else {
		internalSignature = o.InternalSignature
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", o.DateCreated)
	if err != nil {
		log.Printf("error parse time and date: %v\n", err)
		return "", err
	}
	_, err = tx.Exec(`INSERT INTO public.orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, internalSignature,
		o.CustomerID, o.DeliveryService, o.Shardkey, o.SmID, t, o.OofShard,
	)
	if err != nil {
		log.Printf("error insert in order: %v", err)
		tx.Rollback()
		return "", err
	}

	//Items
	for _, item := range o.Items {
		log.Println(item)
		_, err := tx.Exec(`INSERT INTO public.items
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID,
			item.Brand, item.Status,
		)
		if err != nil {
			log.Printf("error insert in items: %v", err)
			tx.Rollback()
			return "", err
		}
	}

	//Payment
	var requestID interface{}
	if o.Payment.RequestID == "" {
		requestID = nil
	} else {
		requestID = o.Payment.RequestID
	}
	_, err = tx.Exec(`INSERT INTO public.payment VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		o.Payment.Transaction, requestID, o.Payment.Currency, o.Payment.Provider,
		o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee,
	)
	if err != nil {
		log.Printf("error insert in payment: %v", err)
		tx.Rollback()
		return "", err
	}

	//Delivery
	_, err = tx.Exec(`INSERT INTO public.delivery VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region,
		o.Delivery.Email,
	)
	if err != nil {
		log.Printf("error insert in delivery: %v", err)
		tx.Rollback()
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("error transaction fix: %v\n", err)
		return "", err
	}
	//ch.Store(o.OrderUID, o)
	return o.OrderUID, nil
}
