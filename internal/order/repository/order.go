package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type orderRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func (o *orderRepository) Delete(ctx context.Context, id int) error {
	query := o.SqlBuilder.Delete("orders").Where(sq.Eq{"id": id})

	res, err := query.RunWith(o.Conn).ExecContext(ctx)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		if err != nil {
			return err
		} else {
			return domain.ErrOrderNotFound
		}
	}

	return nil
}

func (o *orderRepository) Get(ctx context.Context, id int) (*domain.Order, error) {
	query := o.SqlBuilder.Select("id", "pet_id", "ship_date", "status", "complete").From("orders")
	query = query.Where(sq.Eq{"id": id})

	row := query.RunWith(o.Conn).QueryRowContext(ctx)
	var order domain.Order
	if err := row.Scan(&order.Id, &order.PetId, &order.ShipDate, &order.Status, &order.Complete); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}

		return nil, err
	}

	return &order, nil
}

func (o *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := o.SqlBuilder.Insert("orders").Columns("pet_id", "ship_date", "status", "complete")
	query = query.Values(order.PetId, order.ShipDate, order.Status, order.Complete)
	query = query.Suffix("RETURNING id")

	row := query.RunWith(o.Conn).QueryRowContext(ctx)
	return row.Scan(&order.Id)
}

func NewOrderRepository(conn *sql.DB) domain.OrderRepository {
	return &orderRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}
