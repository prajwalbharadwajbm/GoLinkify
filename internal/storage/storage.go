package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type Storage interface {
	StoreURL(url string, isPublic bool) string
	GetURL(shortURL string) (string, error)
	GetURLs() ([]*URL, error)
	GetPublicURLs() ([]*URL, error)
	GetPrivateURLs() ([]*URL, error)
	DeleteUrls() (string, error)
}

type URL struct {
	ShortURL string
	LongURL  string
	IsPublic bool
}

type InMemoryStorage struct {
	mu   sync.Mutex
	urls map[string]*URL
}

func NewMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{urls: make(map[string]*URL)}
}
func (s *InMemoryStorage) StoreURL(url string, isPublic bool) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	hashBytes := hash.Sum(nil)
	shortURL := hex.EncodeToString(hashBytes[:3])
	newURL := &URL{
		ShortURL: shortURL,
		LongURL:  url,
		IsPublic: isPublic,
	}
	s.urls[shortURL] = newURL

	return shortURL
}

func (s *InMemoryStorage) GetURLs() ([]*URL, error) {
	var urls []*URL
	for _, u := range s.urls {
		urls = append(urls, u)
	}
	return urls, nil
}
func (s *InMemoryStorage) GetURL(shortURL string) (string, error) {
	url, ok := s.urls[shortURL]
	if !ok {
		return "", fmt.Errorf("URL not found")
	}
	return url.LongURL, nil
}

func (s *InMemoryStorage) GetPublicURLs() ([]*URL, error) {
	var publicURLs []*URL
	for _, u := range s.urls {
		if u.IsPublic {
			publicURLs = append(publicURLs, u)
		}
	}
	return publicURLs, nil
}
func (s *InMemoryStorage) GetPrivateURLs() ([]*URL, error) {
	var privateURLs []*URL
	for _, u := range s.urls {
		if !u.IsPublic {
			privateURLs = append(privateURLs, u)
		}
	}
	return privateURLs, nil
}

func (s *InMemoryStorage) DeleteUrls() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear the map to delete all URLs
	s.urls = make(map[string]*URL)

	return "All URLs deleted", nil
}
