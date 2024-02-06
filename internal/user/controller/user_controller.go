package controller

import (
	"encoding/json"
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

	r.Post("/user", u.Create)
	r.Get("/user/{username}", u.Get)
	r.Put("/user/{username}", u.Update)
	r.Delete("/user/{username}", u.Delete)
	r.Post("/user/createWithList", u.CreateWithList)

	r.Post("/user/login", u.Login)
	r.Get("/user/logout", u.Logout)
}

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
		w.WriteHeader(http.StatusNotFound)
		u.responder.ErrorInternal(w, err)
		return
	}

	u.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "user updated",
		Data:    nil,
	})
}

func (u *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		u.responder.ErrorBadRequest(w, fmt.Errorf("param username is not set"))
		return
	}

	err := u.userUsecase.Delete(r.Context(), username)
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
		Message: "user deleted",
		Data:    nil,
	})
}

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
