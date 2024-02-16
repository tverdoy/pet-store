package repository

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type photoRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func (p *photoRepository) DeleteByPet(ctx context.Context, petId int) error {
	query := p.SqlBuilder.Delete("photos").Where(sq.Eq{"pet_id": petId})
	_, err := query.RunWith(p.Conn).ExecContext(ctx)

	return err
}

func (p *photoRepository) Create(ctx context.Context, photo domain.PhotoDTO) error {
	query := p.SqlBuilder.Insert("photos").Columns("pet_id").Values(photo.PetId)
	query = query.Suffix("RETURNING id")

	row := query.RunWith(p.Conn).QueryRowContext(ctx)
	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	photo.Id = id
	return nil
}

func (p *photoRepository) GetByPet(ctx context.Context, petId int) ([]*domain.PhotoDTO, error) {
	query := p.SqlBuilder.Select("id", "pet_id").From("photos")
	query = query.Where(sq.Eq{"pet_id": petId})

	rows, err := query.RunWith(p.Conn).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	photos := make([]*domain.PhotoDTO, 0)
	for rows.Next() {
		var photo domain.PhotoDTO
		if err := rows.Scan(&photo.Id, &photo.PetId); err != nil {
			return nil, err
		}

		photos = append(photos, &photo)
	}

	return photos, nil
}

func NewPhotoRepository(conn *sql.DB) domain.PhotoRepository {
	return &photoRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}
