package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"wbl0/internal"
	"wbl0/internal/model"
)

type OrderRepo struct {
	db *sqlx.DB
}

type OrderRepository interface {
	CreateOrder(order model.Order) error
	GetOrderByUid(uid string) (model.Order, error)
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func MustConnection(cfg *internal.Config) *sqlx.DB {

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode))
	if err != nil {
		panic(fmt.Errorf("Error connecting to db: %s", err))
		return nil
	}
	return db
}

func (r *OrderRepo) CreateOrder(order model.Order) error {
	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		fmt.Errorf("Error marshalling delivery json: %s", err)
		return err
	}
	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		fmt.Errorf("Error marshalling payment json: %s", err)
		return err
	}
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		fmt.Errorf("Error marshalling items json: %s", err)
		return err
	}
	checkUid, _ := r.GetOrderByUid(order.OrderUid)
	if checkUid.OrderUid == order.OrderUid {
		return fmt.Errorf("Order with uid %s already exists", order.OrderUid)
	}
	tx := r.db.MustBegin()
	tx.MustExec("INSERT INTO orders (order_uid, track_number, entry, delivery, payment,items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
		order.OrderUid, order.TrackNumber, order.Entry, deliveryJSON, paymentJSON, itemsJSON, order.Locale,
		order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) GetOrderByUid(uid string) (model.Order, error) {
	var order model.Order
	rows, err := r.db.Query("SELECT * FROM orders WHERE order_uid = $1", uid)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Delivery, &order.Payment, &order.Items, &order.Locale,
			&order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		if err != nil {
			return model.Order{}, err
		}
		return order, nil
	}
	return model.Order{
		OrderUid: "не найдено",
	}, sql.ErrNoRows
}
