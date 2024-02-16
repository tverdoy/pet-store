package repository

import (
	"context"
	"database/sql"
	"errors"
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
	res, err := query.RunWith(a.Conn).ExecContext(ctx)

	if n, err := res.RowsAffected(); err != nil || n == 0 {
		if err != nil {
			return err
		} else {
			return domain.ErrSessionNotFound
		}
	}

	return err
}

// ExistsSession checks if a session exists.
// Returns (false, nil) if no session is found.
func (a *authRepository) ExistsSession(ctx context.Context, sessionId int) (bool, error) {
	query := a.SqlBuilder.Select("id").From("auth").Where(sq.Eq{"id": sessionId})
	row := query.RunWith(a.Conn).QueryRowContext(ctx)
	var id int
	err := row.Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return err == nil, err
}
