package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/duckvoid/yago-mart/internal/domain/user"
	"github.com/duckvoid/yago-mart/internal/service"
)

type Handler struct {
	svc    *service.AuthService
	logger *slog.Logger
}

func NewAuthHandler(service *service.AuthService, logger *slog.Logger) *Handler {
	return &Handler{svc: service, logger: logger.With(slog.String("handler", "auth"))}
}

func (a *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.logger.Error("failed to decode request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.svc.Register(r.Context(), req.Login, req.Password); err != nil {
		switch {
		case errors.Is(err, user.ErrAlreadyExist):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	token, err := a.svc.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var respBuf bytes.Buffer
	if err := json.NewEncoder(&respBuf).Encode(RegisterResponse{
		Message: fmt.Sprintf("User %s succesfully register and authenticated", req.Login),
		Code:    http.StatusOK,
	}); err != nil {
		a.logger.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// По ТЗ передаем токен в хедере
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBuf.Bytes()); err != nil {
		a.logger.Error("failed to write response", "error", err)
	}
}

func (a *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.logger.Error("failed to decode request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := a.svc.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var respBuf bytes.Buffer
	if err := json.NewEncoder(&respBuf).Encode(LoginResponse{
		Message: fmt.Sprintf("User %s succesfully login", req.Login),
		Code:    http.StatusOK,
		Token:   token,
	}); err != nil {
		a.logger.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// По ТЗ передаем токен в хедере
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBuf.Bytes()); err != nil {
		a.logger.Error("failed to write response", "error", err)
	}
}
