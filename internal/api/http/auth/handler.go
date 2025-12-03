package authapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: service}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.svc.Register(req.Login, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(RegisterResponse{
		Message: fmt.Sprintf("User %s succesfully register", req.Login),
		Code:    http.StatusOK,
	})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {}
