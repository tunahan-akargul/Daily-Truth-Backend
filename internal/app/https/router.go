package https

import (
	"net/http"

	"main-services/internal/app/https/handlers"
	"main-services/internal/app/words"
	"main-services/internal/db"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type collectionGetter interface {
	Collection(string) *mongo.Collection
}

func NewRouter(mongoClient *db.Mongo) http.Handler {
	router := chi.NewRouter()

	// middlewares: CORS, logging, recoverer, etc.
	// r.Use(...)

	// Build Word feature
	wordRepository := words.NewRepository(mongoClient)
	wordService := words.NewService(wordRepository)
	wordHandler := handlers.NewWordHandler(wordService)

	router.Route("/", func(router chi.Router) {
		// r.Use(AuthMiddleware()) // TODO: Add auth middleware when ready
		router.Post("/post-word", wordHandler.Create)
		//router.Get("/get-words", wordHandler.List)      // later
		// r.Get("/{id}", wordH.GetByID)
		// r.Delete("/{id}", wordH.Delete)
	})

	return router
}
