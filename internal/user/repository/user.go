package repository

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"petstore/internal/domain"
)

type userRepository struct {
	Conn       *sql.DB
	SqlBuilder sq.StatementBuilderType
}

func NewUserRepository(conn *sql.DB) domain.UserRepository {
	return &userRepository{Conn: conn, SqlBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := u.SqlBuilder.Insert("users")
	query = query.Columns("username", "first_name", "last_name", "email", "phone", "password", "user_status")
	query = query.Values(user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.Password, user.UserStatus)
	_, err := query.RunWith(u.Conn).ExecContext(ctx)

	return err
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := u.SqlBuilder.Select("id", "username", "first_name", "last_name", "email", "phone", "password", "user_status")
	query = query.From("users").Where(sq.Eq{"username": username})

	row := query.RunWith(u.Conn).QueryRowContext(ctx)
	user := &domain.User{}
	err := row.Scan(&user.Id, &user.Username, &user.FirstName,
		&user.LastName, &user.Email, &user.Phone, &user.Password, &user.UserStatus)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	return user, err
}

func (u *userRepository) GetIdByUsername(ctx context.Context, username string) (int, error) {
	query := u.SqlBuilder.Select("id")
	query = query.From("users").Where(sq.Eq{"username": username})

	row := query.RunWith(u.Conn).QueryRowContext(ctx)
	var userId int
	err := row.Scan(&userId)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrUserNotFound
	}

	return userId, err
}

func (u *userRepository) Update(ctx context.Context, username string, user *domain.User) error {
	query := u.SqlBuilder.Update("users")
	query = query.Set("username", user.Username).
		Set("first_name", user.FirstName).
		Set("last_name", user.LastName).
		Set("email", user.Email).
		Set("phone", user.Phone).
		Set("password", user.Password).
		Set("user_status", user.UserStatus)

	query = query.Where(sq.Eq{"username": username})

	res, err := query.RunWith(u.Conn).ExecContext(ctx)
	isUpdate, _ := res.RowsAffected()
	if isUpdate == 0 {
		return domain.ErrUserNotFound
	}

	return err
}

func (u *userRepository) Delete(ctx context.Context, username string) error {
	query := u.SqlBuilder.Delete("users").Where(sq.Eq{"username": username})

	res, err := query.RunWith(u.Conn).ExecContext(ctx)
	isDelete, _ := res.RowsAffected()
	if isDelete == 0 {
		return domain.ErrUserNotFound
	}

	return err
}
