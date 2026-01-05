package routes

import (
	"github.com/gorilla/mux"
	"github.com/quynhanh/internship-tracker/internal/application"
	"github.com/quynhanh/internship-tracker/internal/auth"
)

// RegisterRoutes registers all routes in the application
func RegisterRoutes(
	router *mux.Router,
	authHandler *auth.Handler,
	appHandler *application.Handler,
	middleware *auth.Middleware,
) {
	RegisterAuthRoutes(router, authHandler, middleware)
	RegisterApplicationRoutes(router, appHandler, middleware)
}
