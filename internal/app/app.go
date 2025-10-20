package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"main-services/internal/db"
	ihttp "main-services/internal/app/https"
)

type Config struct {
	Port      string
	MongoURI  string
	DBName    string
}

func loadConfig() Config {
	c := Config{
		Port:     getenv("PORT", "8080"),
		MongoURI: getenv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getenv("MONGO_DB", "testdb"),
	}
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Run(ctx context.Context) error {
	cfg := loadConfig()

	mc, err := db.Connect(ctx, cfg.MongoURI, cfg.DBName)
	if err != nil {
		return err
	}
	defer mc.Close(context.Background())

	router := ihttp.NewRouter(mc)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down...")
		shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shCtx)
	case err := <-errCh:
		return err
	}
}