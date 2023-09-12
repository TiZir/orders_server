package db

import (
	"database/sql"
	"log"
	"time"
)

type Repository struct {
	db *sql.DB
	//cache map[string]Order
}

func SelectOrderByID(id string) (Order, error) {
	var o Order
	var r Repository
	var err error

	r.db, err = GetDB()
	if err != nil {
		log.Println(err)
		return o, err
	}

	// Сбор данных об Order
	err = r.db.QueryRow(`SELECT order_uid, track_number, entry, locale, ifnull(internal_signature, ""), customer_id, 
	delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1;`, id).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
		&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard,
	)
	if err != nil {
		log.Println(id, o)
		log.Printf("err select in order: %v", err)
		return o, err
	}

	// Сбор данных о Payment
	err = r.db.QueryRow(`SELECT transaction, ifnull(request_id, ""), currency, provider, amount, payment_dt, bank,
	delivery_cost, goods_total, custom_fee FROM payment WHERE transaction = $1;`, o.OrderUID).Scan(
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider,
		&o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)
	if err != nil {
		log.Println(id, o)
		log.Printf("err select in payment: %v", err)
		return o, err
	}

	// Сбор данных о Delivery
	err = r.db.QueryRow(`SELECT name, phone, zip, city, address, region,
	email FROM delivery WHERE order_uid = $1;`, o.OrderUID).Scan(
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip,
		&o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region,
		&o.Delivery.Email,
	)
	if err != nil {
		log.Println(id, o)
		log.Printf("err select in delivery: %v", err)
		return o, err
	}

	// Сбор данных о Items
	rows, err := r.db.Query(`SELECT * FROM items WHERE track_number = $1;`, o.TrackNumber)
	if err != nil {
		log.Println(id, o)
		log.Printf("err select in items: %v", err)
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
			return o, err
		}
		o.Items = append(o.Items, i)
	}
	return o, nil
}

// Сохранение Order в БД
func AddOrder(o Order) (string, error) {
	var r Repository
	var err error

	r.db, err = GetDB()
	if err != nil {
		log.Println(err)
		return o.OrderUID, err
	}

	// Добавление Order
	var internalSignature interface{}
	if o.InternalSignature == "" {
		internalSignature = nil
	} else {
		internalSignature = o.InternalSignature
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", o.DateCreated)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, err = r.db.Exec(`INSERT INTO public.orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, internalSignature,
		o.CustomerID, o.DeliveryService, o.Shardkey, o.SmID, t, o.OofShard,
	)
	if err != nil {
		log.Printf("err insert in order: %v", err)
		return "", err
	}

	// добавление Items
	for _, item := range o.Items {
		log.Println(item)
		_, err := r.db.Exec(`INSERT INTO public.items
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID,
			item.Brand, item.Status,
		)
		if err != nil {
			log.Printf("err insert in items: %v", err)
			return "", err
		}
	}

	// Добавление Payment
	var requestID interface{}
	if o.Payment.RequestID == "" {
		requestID = nil
	} else {
		requestID = o.Payment.RequestID
	}
	_, err = r.db.Exec(`INSERT INTO public.payment VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		o.Payment.Transaction, requestID, o.Payment.Currency, o.Payment.Provider,
		o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee,
	)
	if err != nil {
		log.Printf("err insert in payment: %v", err)
		return "", err
	}

	// Добавление Delivery
	_, err = r.db.Exec(`INSERT INTO public.delivery VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region,
		o.Delivery.Email,
	)
	if err != nil {
		log.Printf("err insert in delivery: %v", err)
		return "", err
	}

	// После успешной записи добавляем в кеш

	return o.OrderUID, nil
}
