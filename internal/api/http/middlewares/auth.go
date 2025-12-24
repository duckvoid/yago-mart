package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/duckvoid/yago-mart/internal/service"
)

type UserCtxKeyType struct{}

var userCtxKey UserCtxKeyType

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := service.AuthToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtxKey, user)))
	})
}

func UserFromCtx(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(userCtxKey).(string)
	return u, ok
}
