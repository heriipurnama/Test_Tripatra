package middleware

import (
	"context"
	"net/http"
	"strings"

	"procurement-app-backend/utils"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// Allow unauthenticated users in
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Validate JWT token
		tokenStr := strings.Replace(header, "Bearer ", "", 1)
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		// Put it in context
		ctx := context.WithValue(r.Context(), userCtxKey, claims)

		// And call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *utils.Claims {
	raw, _ := ctx.Value(userCtxKey).(*utils.Claims)
	return raw
}
