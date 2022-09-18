package shortner

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"

	"github.com/AllStars123/urlshortner/internal/storages"
)

func AddURL(longURL string, data storages.URLStorage) string {
	shortURL := Shorten(longURL)
	data.SetURL(shortURL, longURL)
	return shortURL
}
func GetURL(shortURL string, data storages.URLStorage) (string, error) {
	result, _ := data.GetURL(shortURL)
	if result == "" {
		return "", errors.New("not found")
	}
	return result, nil
}

func Shorten(longUrl string) string {
	hasher := sha1.New()
	hasher.Write([]byte(longUrl))
	urlHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10]
	return urlHash
}
