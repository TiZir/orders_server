package db

import (
	"database/sql"
	"log"
)

type Repository struct {
	db    *sql.DB
	cache map[string]Order
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
	err = r.db.QueryRow(`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, 
	delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE id = ?`, id).Scan(
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
		&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard,
	)
	if err != nil {
		log.Println(err)
		return o, err
	}

	// Сбор данных о Payment
	err = r.db.QueryRow(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank,
	delivery_cost, goods_total, custom_fee FROM payment WHERE transaction = ?`, o.OrderUID).Scan(
		&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider,
		&o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost,
		&o.Payment, &o.Payment.GoodsTotal, &o.Payment.CustomFee,
	)
	if err != nil {
		log.Println(err)
		return o, err
	}

	// Сбор данных о Delivery
	err = r.db.QueryRow(`SELECT name, phone, zip, city, address, region,
	email FROM delivery WHERE order_uid = ?`, o.OrderUID).Scan(
		&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip,
		&o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region,
		&o.Delivery.Email,
	)
	if err != nil {
		log.Println(err)
		return o, err
	}

	// Сбор данных о Items
	rows, err := r.db.Query(`SELECT * FROM items WHERE track_number = ?`, o.TrackNumber)
	if err != nil {
		log.Println(err)
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

	// добавление Items
	for _, item := range o.Items {
		_, err := r.db.Exec(`INSERT INTO items
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID,
			item.Brand, item.Status,
		)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	// Добавление Payment
	_, err = r.db.Exec(`INSERT INTO payment VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Payment.Transaction, o.Payment.RequestID, o.Payment.Currency, o.Payment.Provider,
		o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost,
		o.Payment, o.Payment.GoodsTotal, o.Payment.CustomFee,
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Добавление Delivery
	_, err = r.db.Exec(`INSERT INTO delivery VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		o.OrderUID, o.Delivery.Name, o.Delivery.Phone, o.Delivery.Zip,
		o.Delivery.City, o.Delivery.Address, o.Delivery.Region,
		o.Delivery.Email,
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Добавление Order
	_, err = r.db.Exec(`INSERT INTO orders VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.Shardkey, o.SmID, o.DateCreated, o.OofShard,
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// После успешной записи добавляем в кеш

	return o.OrderUID, nil
}
