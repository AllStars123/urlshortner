package main

import (
	"github.com/AllStars123/urlshortner/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func main() {
	router := gin.Default()
	data := url.Values{}
	router.POST("/post", handlers.CreateShortURL(data))
	router.GET("/:id", handlers.RetriveShortURL(data))
	http.ListenAndServe(":8080", router)
}
