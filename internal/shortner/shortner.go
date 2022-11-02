package shortner

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func Shorten(longURL string) string {
	splitURL := strings.Split(longURL, "://")
	hasher := sha1.New()
	if len(splitURL) < 2 {
		hasher.Write([]byte(longURL))
	} else {
		hasher.Write([]byte(splitURL[1]))
	}
	urlHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return string(urlHash)
}

//package shortner
//
//import (
//"crypto/sha1"
//"encoding/base64"
//)
//
//func Shorten(longURL string) string {
//	hasher := sha1.New()
//	hasher.Write([]byte(longURL))
//	urlHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10]
//	return urlHash
//}
