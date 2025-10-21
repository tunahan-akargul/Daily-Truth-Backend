package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	ihttp "main-services/internal/app/https"
	"main-services/internal/db"
)

type Config struct {
	Port     string
	MongoURI string
	DBName   string
}

func loadConfig() Config {
	thiDayConfig := Config{
		Port:     getEnvironment("PORT", "8083"),
		MongoURI: getEnvironment("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getEnvironment("MONGO_DB", "main-services"),
	}
	return thiDayConfig
}

func getEnvironment(name, defaultValue string) string {
	if variable := os.Getenv(name); variable != "" {
		return variable
	}
	return defaultValue
}

func Run(ctx context.Context) error {
	config := loadConfig()

	mongoClient, err := db.Connect(ctx, config.MongoURI, config.DBName)
	if err != nil {
		return err
	}
	defer mongoClient.Close(context.Background())

	router := ihttp.NewRouter(mongoClient)

	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on :%s", config.Port)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down...")
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutCtx)
	case err := <-errCh:
		return err
	}
}
