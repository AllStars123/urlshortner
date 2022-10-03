package main

import (
	"net/http"

	"github.com/AllStars123/urlshortner/internal/handlers"
	"github.com/AllStars123/urlshortner/internal/storages"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	data := storages.URLStorage{}
	router.POST("/", handlers.CreateShortURL(data))
	router.GET("/:id", handlers.RetriveShortURL(data))
	http.ListenAndServe(":8080", router)
}
