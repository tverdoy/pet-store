package repository

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type authRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func NewAuthRepository(conn *sql.DB) domain.AuthRepository {
	return &authRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (a *authRepository) RegisterSession(ctx context.Context, userId int) (int, error) {
	query := a.SqlBuilder.Insert("auth")
	query = query.Columns("user_id").Values(userId).Suffix("RETURNING id")

	raw := query.RunWith(a.Conn).QueryRowContext(ctx)
	var sessionId int
	err := raw.Scan(&sessionId)

	return sessionId, err
}

func (a *authRepository) UnregisterSession(ctx context.Context, sessionId int) error {
	query := a.SqlBuilder.Delete("auth").Where(sq.Eq{"id": sessionId})
	_, err := query.RunWith(a.Conn).ExecContext(ctx)

	return err
}

func (a *authRepository) UnregisterAllSession(ctx context.Context, userId int) error {
	query := a.SqlBuilder.Delete("auth").Where(sq.Eq{"user_id": userId})
	_, err := query.RunWith(a.Conn).ExecContext(ctx)

	return err
}

func (a *authRepository) ExistsSession(ctx context.Context, sessionId int) (bool, error) {
	query := a.SqlBuilder.Select("id").From("auth").Where(sq.Eq{"id": sessionId})
	row := query.RunWith(a.Conn).QueryRowContext(ctx)
	var id int
	err := row.Scan(&id)

	return err == nil, err
}
