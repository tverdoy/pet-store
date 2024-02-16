package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type categoryRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func (c *categoryRepository) GetElseCreate(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	categoryDb, err := c.GetByName(ctx, category.Name)

	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			// category not exist, create it
			if err := c.Create(ctx, category); err != nil {
				return nil, err
			}

			return category, nil
		} else {
			// error in db
			return nil, err
		}
	}

	return categoryDb, nil
}

func (c *categoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	query := c.SqlBuilder.Select("id", "name").From("categories").Where(sq.Eq{"name": name})
	row := query.RunWith(c.Conn).QueryRowContext(ctx)

	var category domain.Category
	err := row.Scan(&category.Id, &category.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCategoryNotFound
	}

	return &category, err
}

func NewCategoryRepository(conn *sql.DB) domain.CategoryRepository {
	return &categoryRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (c *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := c.SqlBuilder.Insert("categories").Columns("name").Values(category.Name)
	query = query.Suffix("RETURNING id")
	row := query.RunWith(c.Conn).QueryRowContext(ctx)

	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	category.Id = id

	return nil
}

func (c *categoryRepository) Get(ctx context.Context, id int) (*domain.Category, error) {
	query := c.SqlBuilder.Select("id", "name").From("categories").Where(sq.Eq{"id": id})
	row := query.RunWith(c.Conn).QueryRowContext(ctx)

	var category domain.Category
	err := row.Scan(&category.Id, &category.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrCategoryNotFound
	}

	return &category, err
}
