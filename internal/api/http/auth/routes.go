package authapi

import (
	"github.com/go-chi/chi/v5"
)

func NewAuthRoute(r chi.Router, handler *AuthHandler) {
	r.Route("/", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}
