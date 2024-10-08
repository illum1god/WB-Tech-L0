package postgres

import (
	"WB-L0/internal/configs"
	"WB-L0/internal/domain"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	ordersTable     = "orders"
	deliveriesTable = "deliveries"
	paymentsTable   = "payments"
	itemsTable      = "items"
)

type Order interface {
	GetOrders(ctx context.Context) ([]domain.Order, error)
	GetOrderByUID(ctx context.Context, orderUID string) (domain.Order, error)
	SaveOrder(ctx context.Context, order domain.Order) error
}

func NewPostgresDB(cfg configs.DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode, cfg.Password)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	log.Printf("Подключение к DB: host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

var (
	queryGetOrders     = fmt.Sprintf(`SELECT * FROM %s`, ordersTable)
	queryGetOrderByUID = fmt.Sprintf(`SELECT * FROM %s WHERE uid = $1`, ordersTable)
	queryAddOrder      = fmt.Sprintf(`
        INSERT INTO %s (
            uid, track_number, entry, locale, internal_signature, 
            customer_id, delivery_service, shard_key, sm_id, 
            date_created, oof_shard
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        )
    `, ordersTable)
)

var (
	queryGetDelivery = fmt.Sprintf(`SELECT * FROM %s WHERE order_uid = $1`, deliveriesTable)
	queryAddDelivery = fmt.Sprintf(`
        INSERT INTO %s (
            order_uid, name, phone, zip, city, address, 
            region, email
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
    `, deliveriesTable)
)

var (
	queryGetPayment = fmt.Sprintf(`SELECT * FROM %s WHERE order_uid = $1`, paymentsTable)
	queryAddPayment = fmt.Sprintf(`
        INSERT INTO %s (
            order_uid, transaction, request_id, currency, provider, 
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        )
    `, paymentsTable)
)

var (
	queryGetItems = fmt.Sprintf(`SELECT * FROM %s WHERE order_uid = $1`, itemsTable)
	queryAddItem  = fmt.Sprintf(`
        INSERT INTO %s (
            order_uid, chrt_id, track_number, price, rid, name, 
            sale, size, total_price, nm_id, brand, status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
        )
    `, itemsTable)
)

type orderRepository struct {
	orderDB *sqlx.DB
}

func NewOrder(db *sqlx.DB) Order {
	return &orderRepository{orderDB: db}
}

func (o *orderRepository) GetOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order

	if err := o.orderDB.SelectContext(ctx, &orders, queryGetOrders); err != nil {
		return nil, err
	}

	for i, order := range orders {
		var err error
		orders[i].Delivery, orders[i].Payment, orders[i].Items, err = o.GetDataByUID(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}
	}

	return orders, nil
}

func (o *orderRepository) GetOrderByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	var order domain.Order

	if err := o.orderDB.GetContext(ctx, &order, queryGetOrderByUID, orderUID); err != nil {
		return domain.Order{}, err
	}

	var err error
	order.Delivery, order.Payment, order.Items, err = o.GetDataByUID(ctx, order.OrderUID)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (o *orderRepository) SaveOrder(ctx context.Context, order domain.Order) error {
	tx, err := o.orderDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryAddOrder,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryAddDelivery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryAddPayment,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryAddItem, order.OrderUID,
			item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (o *orderRepository) GetDataByUID(ctx context.Context, uid string) (delivery domain.Delivery, payment domain.Payment, items []domain.Item, err error) {
	if err = o.orderDB.GetContext(ctx, &delivery, queryGetDelivery, uid); err != nil {
		return
	}
	if err = o.orderDB.GetContext(ctx, &payment, queryGetPayment, uid); err != nil {
		return
	}
	err = o.orderDB.SelectContext(ctx, &items, queryGetItems, uid)
	return
}
