package shortner

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

func AddURL(longURL string, data url.Values) string {
	shortURL := Shorten(longURL)
	data.Set(shortURL, longURL)
	return shortURL
}
func GetURL(shortURL string, data url.Values) (string, error) {
	result := data.Get(shortURL)
	if result == "" {
		return "", errors.New("not found")
	}
	return result, nil
}

func Shorten(longUrl string) string {
	split := strings.Split(longUrl, "://")
	hasher := sha1.New()
	if len(split) < 2 {
		hasher.Write([]byte(longUrl))
	} else {
		hasher.Write([]byte(split[1]))
	}
	urlHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10]
	return urlHash
}
