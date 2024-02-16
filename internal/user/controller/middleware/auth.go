package middleware

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"petstore/internal/domain"
	"petstore/internal/responder"
)

func Authenticator(resp responder.Responder, userUsecase domain.UserUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				resp.ErrorUnauthorized(w, err)
				return
			}

			if token == nil {
				resp.ErrorUnauthorized(w, errors.New("token is nil"))
				return
			}

			if isAuth, err := userUsecase.IsAuthenticated(r.Context(), token); err != nil || !isAuth {
				if err != nil {
					resp.ErrorInternal(w, err)
					return
				} else {
					resp.ErrorUnauthorized(w, errors.New("session was logout"))
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
