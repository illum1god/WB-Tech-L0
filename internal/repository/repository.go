package repository

import (
	"WB-L0/internal/repository/order/cache"
	"WB-L0/internal/repository/order/postgres"
	"github.com/jmoiron/sqlx"
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
