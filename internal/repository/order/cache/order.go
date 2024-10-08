package cache

import (
	"WB-L0/internal/domain"
	"sync"
)

type Order interface {
	Get(uid string) (domain.Order, bool)
	Save(uid string, orderToSave *domain.Order) bool
}

type Cache struct {
	Order
}

func NewCache() Cache {
	return Cache{Order: NewOrder()}
}

type orderRepository struct {
	orderRepoByUID map[string]domain.Order
	mu             sync.RWMutex
}

func NewOrder() Order {
	return &orderRepository{orderRepoByUID: make(map[string]domain.Order)}
}

func (o *orderRepository) Get(uid string) (domain.Order, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, has := o.orderRepoByUID[uid]

	if !has {
		return domain.Order{}, false
	}

	return order, true
}

func (o *orderRepository) Save(uid string, orderToSave *domain.Order) bool {
	if orderToSave == nil {
		return false
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	o.orderRepoByUID[uid] = *orderToSave
	return true
}
