package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/AllStars123/urlshortner/cmd/shortener/configuration"
	"github.com/AllStars123/urlshortner/internal/models"
	"github.com/AllStars123/urlshortner/internal/models/mocks"
)

func setupRouter(repo models.RepositoryInterface, baseURL string) *gin.Engine {
	router := gin.Default()

	handler := New(repo, baseURL)

	router.GET("/:id", handler.RetriveShortURL)
	router.POST("/", handler.CreateShortURL)
	router.POST("/api/shorten", handler.ShortenURL)

	router.HandleMethodNotAllowed = true

	return router
}

// Не понимаю почему тут падают тесты
func TestRetriveShortURL(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		request string
		err     error
		result  string
		want    want
	}{
		{
			name:    "GET without id",
			request: "",
			result:  "",
			err:     errors.New("not found"),
			want: want{
				statusCode:  405,
				response:    `405 method not allowed`,
				contentType: `text/plain`,
			},
		},
		{
			name:    "GET with correct id",
			request: "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			result:  "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			err:     nil,
			want: want{
				statusCode:  307,
				response:    ``,
				contentType: `text/plain; charset=utf-8`,
			},
		},
		{
			name:    "GET with incorrect id",
			request: "12398fv58Wr3hGGIzm2-aH2zA628Ng=",
			result:  "",
			err:     errors.New("not found"),
			want: want{
				statusCode:  404,
				response:    `{"detail":"not found"}`,
				contentType: `application/json; charset=utf-8`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoMock := new(mocks.RepositoryInterface)
			repoMock.On("GetURL", tt.request).Return(tt.result, tt.err)

			cfg := configuration.New()
			router := setupRouter(repoMock, cfg.BaseURL)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/"+tt.request, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Header()["Content-Type"][0], tt.want.contentType)

			assert.Equal(t, tt.want.statusCode, w.Code)
			resBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if w.Header()["Content-Type"][0] == `application/json; charset=utf-8` {
				assert.JSONEq(t, tt.want.response, string(resBody))
			} else {
				assert.Equal(t, tt.want.response, string(resBody))
			}

		})
	}
}

func TestCreateShortURL(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		request string
		body    string
		result  string
		want    want
	}{
		{
			name:    "correct POST",
			request: "",
			body:    "https://www.ilovepdf.com/ru",
			result:  "98fv58Wr3hGGIzm2-aH2zA628Ng=",
			want: want{
				statusCode:  201,
				response:    `http://localhost:8080/98fv58Wr3hGGIzm2-aH2zA628Ng=`,
				contentType: `text/plain; charset=utf-8`,
			},
		},
		{
			name:    "incorrect POST",
			request: "123",
			body:    "https://www.ilovepdf.com/ru",
			result:  "",
			want: want{
				statusCode:  405,
				response:    `405 method not allowed`,
				contentType: `text/plain`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoMock := new(mocks.RepositoryInterface)
			repoMock.On("AddURL", tt.body).Return(tt.result, nil)

			cfg := configuration.New()
			router := setupRouter(repoMock, cfg.BaseURL)

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
			if w.Header()["Content-Type"][0] == `application/json; charset=utf-8` {
				assert.JSONEq(t, tt.want.response, string(resBody))
			} else {
				assert.Equal(t, tt.want.response, string(resBody))
			}

		})
	}
}

func TestShortenURL(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		request string
		body    string
		rawData string
		result  string
		want    want
	}{
		{
			name:    "correct POST",
			request: "api/shorten",
			body:    `{"url": "https://www.ilovepdf.com/ru"}`,
			rawData: "https://www.ilovepdf.com/ru",
			result:  "98fv58Wr3hGGIzm2-aH2zA628Ng=",
			want: want{
				statusCode:  201,
				response:    `{"result": "http://localhost:8080/98fv58Wr3hGGIzm2-aH2zA628Ng="}`,
				contentType: `application/json; charset=utf-8`,
			},
		},
		{
			name:    "incorrect POST",
			request: "api/shorten",
			body:    `{"url2": "https://www.ilovepdf.com/ru"}`,
			rawData: "https://www.ilovepdf.com/ru",
			result:  "98fv58Wr3hGGIzm2-aH2zA628Ng=",
			want: want{
				statusCode:  400,
				response:    `{"detail": "Bad request"}`,
				contentType: `application/json; charset=utf-8`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoMock := new(mocks.RepositoryInterface)
			repoMock.On("AddURL", tt.rawData).Return(tt.result, nil)

			cfg := configuration.New()
			router := setupRouter(repoMock, cfg.BaseURL)

			body := strings.NewReader(tt.body)
			w := httptest.NewRecorder()
			fmt.Println(tt.request)
			req, _ := http.NewRequest(http.MethodPost, "/"+tt.request, body)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.contentType, w.Header()["Content-Type"][0])

			assert.Equal(t, tt.want.statusCode, w.Code)
			resBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if w.Header()["Content-Type"][0] == `application/json; charset=utf-8` {
				assert.JSONEq(t, tt.want.response, string(resBody))
			} else {
				assert.Equal(t, tt.want.response, string(resBody))
			}

		})
	}
}
