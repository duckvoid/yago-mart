package auth

import (
	"github.com/go-chi/chi/v5"
)

func NewAuthRoute(r chi.Router, handler *Handler) {
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
}
