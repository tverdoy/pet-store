package domain

import (
	"context"
	"errors"
	"time"
)

type OrderStatus string

var ErrOrderNotFound = errors.New("order not found")

const (
	PlacedOrderStatus    OrderStatus = "placed"
	ApprovedOrderStatus  OrderStatus = "approved"
	DeliveredOrderStatus OrderStatus = "delivered"
)

type Order struct {
	Id       int         `json:"id"`
	PetId    int         `json:"petId"`
	ShipDate time.Time   `json:"shipDate"`
	Status   OrderStatus `json:"status"`
	Complete bool        `json:"complete"`
}

type OrderUsecase interface {
	Get(ctx context.Context, id int) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id int) error
}

type OrderRepository interface {
	Get(ctx context.Context, id int) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id int) error
}
