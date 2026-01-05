package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quynhanh/internship-tracker/internal/auth"
)

// RegisterAuthRoutes registers authentication-related routes
func RegisterAuthRoutes(router *mux.Router, authHandler *auth.Handler, m *auth.Middleware) {
	router.HandleFunc("/intern/auth/register", authHandler.Register)
	router.HandleFunc("/intern/auth/login", authHandler.Login)
	router.Handle(
		"/intern/auth/logout",
		m.AuthMiddleware(http.HandlerFunc(authHandler.Logout)),
	)
}
