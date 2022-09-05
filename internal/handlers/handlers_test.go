package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupRuter(data url.Values) *gin.Engine {
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
		longUrl  string
		shortUrl string
		want     want
	}{
		{
			name:     "Get with incorrect id",
			request:  "65463fv58Wr3hdskjIzm2*nmH2zA628Bg=",
			longUrl:  "https://www.ilovepdf.com/ru",
			shortUrl: "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  404,
				response:    `{"detail":"not found"}`,
			},
		},
		{
			name:     "Get with correct id",
			request:  "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			longUrl:  "https://www.ilovepdf.com/ru",
			shortUrl: "tN1fptfy6FNYOpaVrMtQuusk1Po=",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  307,
				response:    "",
			},
		},
	}
	data := url.Values{}
	for _, tt := range tests {
		data.Set(tt.shortUrl, tt.longUrl)
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
				statusCode:  201,
				response:    "http://localhost:8080/t4oQqVMwG34xz4LLJELnjySD7aM=",
			},
		},
		{
			name:    "Incorrect Post",
			request: "test",
			want: want{
				contentType: "text/plain",
				statusCode:  404,
				response:    "404 page not found",
			},
		},
	}
	data := url.Values{}
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
