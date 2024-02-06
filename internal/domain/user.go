package domain

import "context"

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	UserStatus int    `json:"userStatus"`
}

type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, username string, user *User) error
	Delete(ctx context.Context, username string) error
	CreateList(ctx context.Context, users []*User) error
	Login(ctx context.Context, username string, password string) (string, error)
	Logout(ctx context.Context, token string) error
	IsAuthenticated(ctx context.Context, token string) (bool, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, username string, user *User) error
	Delete(ctx context.Context, username string) error
	GetIdByUsername(ctx context.Context, username string) (int, error)
}

type AuthRepository interface {
	RegisterSession(ctx context.Context, userId int) (int, error)
	UnregisterSession(ctx context.Context, sessionId int) error
	UnregisterAllSession(ctx context.Context, userId int) error
	ExistsSession(ctx context.Context, sessionId int) (bool, error)
}
