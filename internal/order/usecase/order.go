package usecase

import (
	"context"
	"petstore/internal/domain"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
}

func (o *orderUsecase) Get(ctx context.Context, id int) (*domain.Order, error) {
	return o.orderRepo.Get(ctx, id)
}

func (o *orderUsecase) Create(ctx context.Context, order *domain.Order) error {
	return o.orderRepo.Create(ctx, order)
}

func (o *orderUsecase) Delete(ctx context.Context, id int) error {
	return o.orderRepo.Delete(ctx, id)
}

func NewOrderUsecase(orderRepo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo}
}
