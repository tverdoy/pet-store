package usecase

import (
	"context"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"petstore/internal/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
	authRepo domain.AuthRepository
	jwtAuth  *jwtauth.JWTAuth
}

func NewUserUsecase(ur domain.UserRepository, au domain.AuthRepository, jwt *jwtauth.JWTAuth) domain.UserUsecase {
	return &userUsecase{userRepo: ur, authRepo: au, jwtAuth: jwt}
}

func (u *userUsecase) hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash)
}

// checkPassword - compare hashed password with target password.
// Return true if password is correct.
func (u *userUsecase) checkPassword(hashedPassword string, targetPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(targetPassword))

	return err == nil
}

func (u *userUsecase) Create(ctx context.Context, user *domain.User) error {
	user.Password = u.hashPassword(user.Password)

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Get(ctx context.Context, username string) (*domain.User, error) {
	return u.userRepo.GetByUsername(ctx, username)
}

func (u *userUsecase) Update(ctx context.Context, username string, user *domain.User) error {
	user.Password = u.hashPassword(user.Password)

	return u.userRepo.Update(ctx, username, user)
}

// Delete - delete user by username and delete all session of this user
func (u *userUsecase) Delete(ctx context.Context, username string) error {
	userId, err := u.userRepo.GetIdByUsername(ctx, username)
	if err != nil {
		return err
	}

	err = u.authRepo.UnregisterAllSession(ctx, userId)
	if err != nil {
		return err
	}

	return u.userRepo.Delete(ctx, username)
}

func (u *userUsecase) CreateList(ctx context.Context, users []*domain.User) error {
	for i, user := range users {
		err := u.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("failed create user %d: %w", i, err)
		}
	}

	return nil
}

func (u *userUsecase) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := u.Get(ctx, username)
	if err != nil {
		return "", err
	}

	if !u.checkPassword(user.Password, password) {
		return "", fmt.Errorf("wrong password")
	}

	sessionId, err := u.authRepo.RegisterSession(ctx, user.Id)
	if err != nil {
		return "", err
	}

	_, token, _ := u.jwtAuth.Encode(map[string]interface{}{"session_id": sessionId})
	return token, nil
}

func (u *userUsecase) Logout(ctx context.Context, token string) error {
	decodeToken, err := u.jwtAuth.Decode(token)
	if err != nil {
		return err
	}

	sessionId, ok := decodeToken.Get("session_id")
	if !ok {
		return fmt.Errorf("session_id not found")
	}

	return u.authRepo.UnregisterSession(ctx, int(sessionId.(float64)))
}

func (u *userUsecase) IsAuthenticated(ctx context.Context, token string) (bool, error) {
	decodeToken, err := u.jwtAuth.Decode(token)
	if err != nil {
		return false, err
	}

	sessionId, ok := decodeToken.Get("session_id")
	if !ok {
		return false, fmt.Errorf("session_id not found")
	}

	return u.authRepo.ExistsSession(ctx, sessionId.(int))
}
