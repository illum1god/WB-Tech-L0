package repository

import (
	"WB-L0/internal/repository/order/cache"
	"WB-L0/internal/repository/order/postgres"
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type Repository struct {
	postgres.Order
	cache.Cache
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		Order: postgres.NewOrder(db),
		Cache: cache.NewCache(),
	}
}

func (r *Repository) RestoreCache() error {
	orders, err := r.Order.GetOrders(context.Background())
	if err != nil {
		return err
	}

	for _, order := range orders {
		err := r.Cache.Save(order.OrderUID, order)
		if err != true {
			log.Printf("Ошибка при добавлении заказа UID %s в кэш: %v", order.OrderUID, err)
			continue
		}
	}

	log.Printf("Кэш успешно восстановлен из базы данных. Загружено %d заказов.", len(orders))
	return nil
}
