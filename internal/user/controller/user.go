package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"petstore/internal/domain"
	"petstore/internal/responder"
)

type UserController struct {
	userUsecase domain.UserUsecase
	responder   responder.Responder
}

func NewUserController(r chi.Router, resp responder.Responder, us domain.UserUsecase) {
	u := UserController{userUsecase: us, responder: resp}
	r.Route("/user", func(r chi.Router) {
		r.Post("/login", u.Login)
		r.Get("/logout", u.Logout)
		r.Post("/createWithList", u.CreateWithList)

		r.Post("/", u.Create)
		r.Get("/{username}", u.Get)
		r.Put("/{username}", u.Update)
		r.Delete("/{username}", u.Delete)
	})
}

// Create this function creates a new user
//
// @Summary		Create a new user
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		pet		body		domain.User			true	"User to add to the store"
// @Success		200		{string}	string				"User created"
// @Failure		400		{string}	string				"Invalid input"
// @Router		/user 	[post]
func (u *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var userInput domain.User
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	if err := u.userUsecase.Create(r.Context(), &userInput); err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "user created",
		Data:    nil,
	})
}

// Get this function is used to get a user
//
// @Summary		Get user by username
// @Tags		user
// @Produce		json
//
// @Param		username path		string				true	"Username of user to return"
//
// @Success		200		{object}	domain.User			"Find user by Username"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"User not found"
// @Router		/user/{username} 		[get]
func (u *UserController) Get(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		u.responder.ErrorBadRequest(w, fmt.Errorf("param username is not set"))
		return
	}

	user, err := u.userUsecase.Get(r.Context(), username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		u.responder.OutputJSON(w, responder.Response{
			Success: false,
			Message: "user not found",
			Data:    nil,
		})

		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "get user",
		Data:    user,
	})
}

// Update this function is used to update a user
//
// @Summary		Update a user with form data
// @Tags		user
// @Accept		json
// @Produce		json
//
// @Param		user	body		domain.User			true	"User object that needs to update"
//
// @Success		200		{string}	string				"User updated"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"User not found"
// @Router		/user 	[put]
func (u *UserController) Update(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		u.responder.ErrorBadRequest(w, fmt.Errorf("param username is not set"))
		return
	}

	var userInput domain.User
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	if err := u.userUsecase.Update(r.Context(), username, &userInput); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			u.responder.ErrorNotFound(w, err)
		} else {
			u.responder.ErrorInternal(w, err)
		}
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "user updated",
		Data:    nil,
	})
}

// Delete this function is used to delete a user.
//
// @Summary		Delete a user by username
// @Tags		user
// @Produce		json
//
// @Param		username path		string				true	"Username of user to delete"
//
// @Success		200		{string}	string				"User deleted"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"User not found"
// @Router		/user/{username} 	[delete]
func (u *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		u.responder.ErrorBadRequest(w, fmt.Errorf("param username is not set"))
		return
	}

	err := u.userUsecase.Delete(r.Context(), username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			u.responder.ErrorNotFound(w, err)
		} else {
			u.responder.ErrorInternal(w, err)
		}
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "user deleted",
		Data:    nil,
	})
}

// CreateWithList this function creates a new users
//
// @Summary		Create a list of new users
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		users		body		[]domain.User		true	"Users to add to the store"
// @Success		200		{string}	string				"Users created"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"User not found"
// @Router		/user/createWithList	[post]
func (u *UserController) CreateWithList(w http.ResponseWriter, r *http.Request) {
	var userInput []*domain.User
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	if err := u.userUsecase.CreateList(r.Context(), userInput); err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "users created",
		Data:    nil,
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login this function login a user
//
// @Summary		Login a user
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		credentials			body				LoginRequest		true	"User credentials"
// @Success		200		{string}	string				"User login"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"User not found"
// @Router		/user/login			[post]
func (u *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginInput LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginInput); err != nil {
		u.responder.ErrorBadRequest(w, err)
		return
	}

	token, err := u.userUsecase.Login(r.Context(), loginInput.Username, loginInput.Password)
	if err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "user login",
		Data:    token,
	})
}

// Logout this function logout a user
//
// @Summary		Logout a user
// @Tags		user
// @Security 	ApiKeyAuth
// @Accept		json
// @Produce		json
//
// @Success		200		{string}	string				"User logout"
//
// @Failure		400		{string}	string				"Invalid input"
// @Router		/user/logout		[get]
func (u *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	token := jwtauth.TokenFromHeader(r)
	if token == "" {
		u.responder.ErrorBadRequest(w, fmt.Errorf("token is not set"))
		return
	}

	err := u.userUsecase.Logout(r.Context(), token)
	if err != nil {
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "user logout",
		Data:    nil,
	})
}
