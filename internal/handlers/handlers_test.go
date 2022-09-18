package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AllStars123/urlshortner/internal/storages"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRuter(data storages.URLStorage) *gin.Engine {
	r := gin.Default()
	r.GET("/:id", RetriveShortURL(data))
	r.POST("/", CreateShortURL(data))
	return r
}

func TestRetriveShortURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name     string
		request  string
		longURL  string
		shortURL string
		want     want
	}{
		{
			name:     "Get with incorrect id",
			request:  "65463fv58Wr3hdskjIzm2*nmH2zA628Bg=",
			longURL:  "https://www.ilovepdf.com/ru",
			shortURL: "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  http.StatusNotFound,
				response:    `{"errorDetails":"not found"}`,
			},
		},
		{
			name:     "Get with correct id",
			request:  "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			longURL:  "https://www.ilovepdf.com/ru",
			shortURL: "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusTemporaryRedirect,
				response:    "",
			},
		},
	}
	data := storages.URLStorage{}
	for _, tt := range tests {
		data.SetURL(tt.shortURL, tt.longURL)
		router := setupRuter(data)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+tt.request, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, w.Header()["Content-Type"][0], tt.want.contentType)

		assert.Equal(t, tt.want.statusCode, w.Code)
		resBody, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatal(err)
		}
		if w.Header()["Content-Type"][0] == "application/json; charset=utf-8" {
			assert.JSONEq(t, tt.want.response, string(resBody))
		} else {
			assert.Equal(t, tt.want.response, string(resBody))
		}
	}

}

func TestCreateShortURL(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "Correct Post",
			request: "",
			body:    "https://www.ilovepdf.com/ru",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/qyfwSuz6En",
			},
		},
		{
			name:    "Incorrect Post",
			request: "test",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusNotFound,
				response:    "404 page not found",
			},
		},
		{
			name:    "Incorrect URL",
			request: "",
			body:    "kl.www.ilovepdf.com/ru",
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    `{"error": "parse \"kl.www.ilovepdf.com/ru\": invalid URI for request"}`,
			},
		},
	}
	data := storages.URLStorage{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRuter(data)
			body := strings.NewReader(tt.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/"+tt.request, body)
			router.ServeHTTP(w, req)
			assert.Equal(t, w.Header()["Content-Type"][0], tt.want.contentType)
			assert.Equal(t, tt.want.statusCode, w.Code)
			resBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if w.Header()["Content-Type"][0] == "application/json; charset=utf-8" {
				assert.JSONEq(t, tt.want.response, string(resBody))
			} else {
				assert.Equal(t, tt.want.response, string(resBody))
			}
		})
	}
}
