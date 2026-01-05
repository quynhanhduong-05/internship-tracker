package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quynhanh/internship-tracker/internal/application"
	"github.com/quynhanh/internship-tracker/internal/auth"
)

// RegisterApplicationRoutes registers application-related routes
func RegisterApplicationRoutes(
	router *mux.Router,
	appHandler *application.Handler,
	m *auth.Middleware,
) {
	router.Handle(
		"/applications",
		m.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				appHandler.Create(w, r)
			case http.MethodGet:
				appHandler.List(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods(http.MethodPost, http.MethodGet)

	router.Handle(
		"/application/{applicationID:[0-9]+}",
		m.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPut:
				appHandler.Update(w, r)
			case http.MethodDelete:
				appHandler.Delete(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	).Methods(http.MethodPut, http.MethodDelete)
}
