package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"main-services/internal/app/https/handlers"
	"main-services/internal/db"
)

func NewRouter(m *db.Mongo) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * 1e9)) // 30s

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })

	// Auth/user endpoints
	r.Route("/", func(r chi.Router) {
		r.Get("/check-email", handlers.CheckEmail(m))
		r.Post("/signup", handlers.SignUp(m))
		r.Post("/signin", handlers.SignIn(m))
		r.Get("/get-details", handlers.GetUserDetails(m))
	})

	return r
}
