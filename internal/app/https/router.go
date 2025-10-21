package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"main-services/internal/app/https/handlers"
	"main-services/internal/db"
)

func NewRouter(myMongo *db.Mongo) http.Handler {
	route := chi.NewRouter()
	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(middleware.Timeout(30 * 1e9)) // 30s

	route.Get("/health", func(wrıter http.ResponseWriter, request *http.Request) { wrıter.Write([]byte("ok")) })

	// Auth/user endpoints
	route.Route("/", func(route chi.Router) {
		route.Get("/check-email", handlers.CheckEmail(myMongo))
		route.Post("/signup", handlers.SignUp(myMongo))
		route.Post("/signin", handlers.SignIn(myMongo))
		route.Get("/get-details", handlers.GetUserDetails(myMongo))
	})

	return route
}
