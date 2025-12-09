package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/domain/user"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *Handler {
	return &Handler{svc: service}
}

func (a *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.svc.Register(req.Login, req.Password); err != nil {
		switch {
		case errors.Is(err, user.ErrAlreadyExist):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(RegisterResponse{
		Message: fmt.Sprintf("User %s succesfully register", req.Login),
		Code:    http.StatusOK,
	})
}

func (a *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := a.svc.Login(req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{
		Message: fmt.Sprintf("User %s succesfully login", req.Login),
		Code:    http.StatusOK,
		Token:   token,
	})
}
