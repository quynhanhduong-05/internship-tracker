package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/quynhanh/internship-tracker/internal/model"
	"gorm.io/gorm"
)

// to avoid key collisions when storing values in context
type contextKey string

// context key used to store the authenticated user's ID
const userIDKey contextKey = "userID"

// Hold dependencies required by authentication middleware
// such as database access for token revocation checks
type Middleware struct {
	DB *gorm.DB
}

func NewMiddleware(db *gorm.DB) *Middleware {
	return &Middleware{DB: db}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token string
		tokenStr := GetTokenFromHeader(r)
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Extract user ID from token claims
		claims, err := ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		userID := uint(claims.UserID)

		// Check revoked token
		var revoked model.RevokedToken
		err = m.DB.Where("token = ?", tokenStr).First(&revoked).Error
		if err == nil {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		} else if err != gorm.ErrRecordNotFound {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if token expires
		if time.Now().After(claims.ExpiresAt.Time) {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		// Store authenticated user ID in request context
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) uint {
	return ctx.Value(userIDKey).(uint)
}
