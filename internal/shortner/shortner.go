package shortner

import (
	"crypto/sha1"
	"encoding/base64"
)

func Shorten(longURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(longURL))
	urlHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10]
	return urlHash
}
