package main

import (
	"net/http"
	"net/url"

	"github.com/AllStars123/urlshortner/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	data := url.Values{}
	router.POST("/", handlers.CreateShortURL(data))
	router.GET("/:id", handlers.RetriveShortURL(data))
	http.ListenAndServe(":8080", router)
}
