package handlers

import (
	"github.com/AllStars123/urlshortner/internal/shortner"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RetriveShortURL(data url.Values) func(c *gin.Context) {
	return func(c *gin.Context) {
		result := map[string]string{}
		long, err := shortner.GetURL(c.Param("id"), data)
		if err != nil {
			result["detail"] = err.Error()
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
			result["detail"] = "Bad request"
			c.IndentedJSON(http.StatusBadRequest, result)
			return
		}
		short := shortner.AddURL(string(body), data)
		c.String(http.StatusCreated, "http://localhost:8080/"+short)
	}
}
