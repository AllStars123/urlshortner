package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/AllStars123/urlshortner/internal/shortner"
	"github.com/gin-gonic/gin"
)

func RetriveShortURL(data url.Values) func(c *gin.Context) {
	return func(c *gin.Context) {
		result := map[string]string{}
		long, err := shortner.GetURL(c.Param("id"), data)
		if err != nil {
			result["errorDetails"] = err.Error()
			c.IndentedJSON(http.StatusNotFound, result)
			return
		}
		c.Header("Location", long)
		c.String(http.StatusTemporaryRedirect, "")
	}
}

func CreateShortURL(data url.Values) func(c *gin.Context) {
	return func(c *gin.Context) {
		result := map[string]string{}
		defer c.Request.Body.Close()
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, result)
			return
		}
		_, err = url.ParseRequestURI(string(body))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		short := shortner.AddURL(string(body), data)
		c.String(http.StatusCreated, "http://localhost:8080/"+short)
	}
}
