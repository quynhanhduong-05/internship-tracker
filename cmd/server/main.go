package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quynhanh/internship-tracker/internal/application"
	"github.com/quynhanh/internship-tracker/internal/auth"
	"github.com/quynhanh/internship-tracker/internal/db"
	"github.com/quynhanh/internship-tracker/internal/model"
	"github.com/quynhanh/internship-tracker/internal/routes"
)

func main() {
	database := db.ConnectMySQL()

	database.AutoMigrate(
		&model.User{},
		&model.Application{},
		&model.RevokedToken{},
	)

	authService := auth.NewService(database)
	authHandler := auth.NewHandler(authService)

	appService := application.NewService(database)
	appHandler := application.NewHandler(appService)
	middleware := auth.NewMiddleware(database)

	r := mux.NewRouter()
	routes.RegisterRoutes(r, authHandler, appHandler, middleware)

	log.Println("Server running on :8080")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
