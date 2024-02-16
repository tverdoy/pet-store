package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type tagRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func (t *tagRepository) AddTagsToPet(ctx context.Context, petId int, tagIds []int) error {
	query := t.SqlBuilder.Insert("pets_tags").Columns("pet_id", "tag_id")

	for _, tagId := range tagIds {
		query = query.Values(petId, tagId)
	}

	_, err := query.RunWith(t.Conn).ExecContext(ctx)

	return err
}

func (t *tagRepository) GetPetTags(ctx context.Context, petId int) ([]*domain.Tag, error) {
	query := t.SqlBuilder.Select("tags.id", "tags.name").From("pets_tags")
	query = query.Join("tags ON tags.id = pets_tags.tag_id")
	query = query.Where(sq.Eq{"pets_tags.pet_id": petId})

	rows, err := query.RunWith(t.Conn).QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tags := make([]*domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			return nil, err
		}

		tags = append(tags, &tag)
	}

	return tags, nil
}

func (t *tagRepository) RemovePetTags(ctx context.Context, petId int) error {
	query := t.SqlBuilder.Delete("pets_tags").Where(sq.Eq{"pet_id": petId})
	_, err := query.RunWith(t.Conn).ExecContext(ctx)

	return err
}

func (t *tagRepository) GetElseCreate(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	categoryDb, err := t.GetByName(ctx, tag.Name)

	if err != nil {
		if errors.Is(err, domain.ErrTagNotFound) {
			// category not exist, create it
			if err := t.Create(ctx, tag); err != nil {
				return nil, err
			}

			return tag, nil
		} else {
			// error in db
			return nil, err
		}
	}

	return categoryDb, nil
}

func (t *tagRepository) GetByName(ctx context.Context, name string) (*domain.Tag, error) {
	query := t.SqlBuilder.Select("id", "name").From("tags").Where(sq.Eq{"name": name})
	row := query.RunWith(t.Conn).QueryRowContext(ctx)

	var tag domain.Tag
	err := row.Scan(&tag.Id, &tag.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrTagNotFound
	}

	return &tag, err
}

func (t *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	query := t.SqlBuilder.Insert("tags").Columns("name").Values(tag.Name)
	query = query.Suffix("RETURNING id")
	row := query.RunWith(t.Conn).QueryRowContext(ctx)

	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	tag.Id = id

	return nil
}

func NewTagRepository(conn *sql.DB) domain.TagRepository {
	return &tagRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}
