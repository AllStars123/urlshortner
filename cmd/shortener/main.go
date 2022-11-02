package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/AllStars123/urlshortner/cmd/shortener/configuration"
	"github.com/AllStars123/urlshortner/internal/handlers"
	"github.com/AllStars123/urlshortner/internal/models"
	"github.com/gin-gonic/gin"
)

func setupRouter(repo models.RepositoryInterface, baseURL string) *gin.Engine {
	router := gin.Default()

	handler := handlers.New(repo, baseURL)

	router.GET("/:id", handler.RetriveShortURL)
	router.POST("/", handler.CreateShortURL)
	router.POST("/api/shorten", handler.ShortenURL)

	router.HandleMethodNotAllowed = true

	return router
}

func main() {

	cfg := configuration.New()

	handler := setupRouter(models.NewFileRepository(cfg.FilePath), cfg.BaseURL)

	server := &http.Server{
		Addr:    cfg.ServerAdress,
		Handler: handler,
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Fatal(server.ListenAndServe())
		cancel()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}
	server.Shutdown(ctx)
}
