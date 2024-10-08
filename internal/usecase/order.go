package usecase

import (
	"WB-L0/internal/domain"
	"WB-L0/internal/repository"
	"context"
)

type OrderInput struct {
	Order domain.Order `json:"order"`
}

type Order interface {
	GetOrders(ctx context.Context) ([]domain.Order, error)
	GetOrderByUID(ctx context.Context, orderUID string) (domain.Order, error)
	SaveOrder(ctx context.Context, orderToSave domain.Order) error
}

type Service struct {
	Order
}

func NewService(repos repository.Repository) Service {
	return Service{Order: NewOrder(repos)}
}

type orderRepository struct {
	repo repository.Repository
}

func NewOrder(repo repository.Repository) Order {
	return &orderRepository{repo: repo}
}

func (o *orderRepository) GetOrders(ctx context.Context) ([]domain.Order, error) {
	return o.repo.GetOrders(ctx)
}

func (o *orderRepository) GetOrderByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	orderRes, ok := o.repo.Get(orderUID)
	if ok {
		return orderRes, nil
	}
	orderRes, err := o.repo.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return domain.Order{}, err
	}
	o.repo.Save(orderUID, orderRes)
	return orderRes, nil
}

func (o *orderRepository) SaveOrder(ctx context.Context, order domain.Order) error {
	err := o.repo.SaveOrder(ctx, order)
	if err != nil {
		return err
	}
	o.repo.Save(order.OrderUID, order)
	return nil
}
