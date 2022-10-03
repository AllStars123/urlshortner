package storages

type URLStorage map[string]string

func (us URLStorage) SetURL(shortURL string, originalURL string) {
	us[shortURL] = originalURL
}

func (us URLStorage) GetURL(shortURL string) (string, bool) {
	originalURL, ok := us[shortURL]
	if !ok {
		return "", false
	}

	return originalURL, true
}
