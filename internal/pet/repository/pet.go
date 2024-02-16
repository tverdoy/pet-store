package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type petRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func (p *petRepository) GetByStatus(ctx context.Context, status domain.PetStatus) ([]*domain.PetDTO, error) {
	query := p.SqlBuilder.Select("id", "category_id", "name", "status")
	query = query.From("pets").Where(sq.Eq{"status": status})

	rows, err := query.RunWith(p.Conn).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pets := make([]*domain.PetDTO, 0)
	for rows.Next() {
		var pet domain.PetDTO
		if err := rows.Scan(&pet.Id, &pet.CategoryId, &pet.Name, &pet.Status); err != nil {
			return nil, err
		}

		pets = append(pets, &pet)
	}

	return pets, err
}

func (p *petRepository) Get(ctx context.Context, id int) (*domain.PetDTO, error) {
	query := p.SqlBuilder.Select("id", "category_id", "name", "status")
	query = query.From("pets").Where(sq.Eq{"id": id})

	row := query.RunWith(p.Conn).QueryRowContext(ctx)
	var pet domain.PetDTO
	err := row.Scan(&pet.Id, &pet.CategoryId, &pet.Name, &pet.Status)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPetNotFound
	}

	return &pet, err
}

func (p *petRepository) Update(ctx context.Context, pet *domain.PetDTO) error {
	query := p.SqlBuilder.Update("pets")
	query = query.Set("category_id", pet.CategoryId).Set("name", pet.Name).Set("status", pet.Status)
	query = query.Where(sq.Eq{"id": pet.Id})

	res, err := query.RunWith(p.Conn).ExecContext(ctx)

	isUpdate, _ := res.RowsAffected()
	if isUpdate == 0 {
		return domain.ErrPetNotFound
	}

	return err
}

func (p *petRepository) Delete(ctx context.Context, id int) error {
	query := p.SqlBuilder.Delete("pets").Where(sq.Eq{"id": id})
	res, err := query.RunWith(p.Conn).ExecContext(ctx)

	isDelete, _ := res.RowsAffected()
	if isDelete == 0 {
		return domain.ErrPetNotFound
	}

	return err
}

func (p *petRepository) Create(ctx context.Context, pet *domain.PetDTO) error {
	query := p.SqlBuilder.Insert("pets").Columns("category_id", "name", "status")
	query = query.Values(pet.CategoryId, pet.Name, pet.Status)
	query = query.Suffix("RETURNING id")

	row := query.RunWith(p.Conn).QueryRowContext(ctx)
	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	pet.Id = id
	return nil
}

func NewPetRepository(conn *sql.DB) domain.PetRepository {
	return &petRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}
