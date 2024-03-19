package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"com.github/prajwalbharadwajbm/GoLinkify/internal/storage"
)

var store storage.Storage

func SetStorage(s storage.Storage) {
	store = s
}

// In the standard library, a handler is an interface that defines the method signature ServeHTTP(w http.ResponseWriter, r *http.Request)
// So, to create a handler, you need to create a struct and implement ServeHTTP
type URLHandler struct{}

func (h *URLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && UrlRe.MatchString(r.URL.Path):
		h.shortenUrl(w, r)
		return
	case r.Method == http.MethodGet && UrlReWithID.MatchString(r.URL.Path):
		h.redirectUrl(w, r)
		return
	case r.Method == http.MethodGet:
		h.getUrls(w, r)
		return
	case r.Method == http.MethodGet && UrlRePublic.MatchString(r.URL.Path):
		h.listPublicShortendUrl(w, r)
		return
	case r.Method == http.MethodGet && UrlRePrivate.MatchString(r.URL.Path):
		h.listPrivateShortendUrl(w, r)
		return
	case r.Method == http.MethodDelete && UrlReWithID.MatchString(r.URL.Path):
		h.deletedShortendUrl(w, r)
		return
	default:
		w.Write([]byte("Home Page"))
	}
}

func (h *URLHandler) shortenUrl(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode JSON request body", http.StatusBadRequest)
		return
	}

	// Extract long URL from JSON data
	longUrl, ok := data["url"]
	if !ok {
		http.Error(w, "URL not found in request body", http.StatusBadRequest)
		return
	}
	shortUrl := store.StoreURL(longUrl, true)
	w.Write([]byte(shortUrl))
}

func (h *URLHandler) redirectUrl(w http.ResponseWriter, r *http.Request) {
	urlString := r.URL.String()

	// Parse the URL string
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %s", err), http.StatusBadRequest)
		return
	}

	// Extract just the path from the parsed URL
	path := parsedURL.Path
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	longUrl, err := store.GetURL(path)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longUrl, http.StatusFound)
}

func (h *URLHandler) getUrls(w http.ResponseWriter, r *http.Request) {
	urls, err := store.GetURLs()
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	responseJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *URLHandler) listPublicShortendUrl(w http.ResponseWriter, r *http.Request) {
	urls, err := store.GetPublicURLs()
	if err != nil {
		http.Error(w, "Failed to retrieve public URLs", http.StatusInternalServerError)
		return
	}
	responseJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *URLHandler) listPrivateShortendUrl(w http.ResponseWriter, r *http.Request) {
	urls, err := store.GetPrivateURLs()
	if err != nil {
		http.Error(w, "Failed to retrieve public URLs", http.StatusInternalServerError)
		return
	}
	responseJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func (h *URLHandler) deletedShortendUrl(w http.ResponseWriter, r *http.Request) {
	response, err := store.DeleteUrls()
	if err != nil {
		http.Error(w, "Failed to delete URLs", http.StatusInternalServerError)
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

var (
	UrlRe = regexp.MustCompile(`^/shorten/*$`)
	// TODO: Improve regex validation to capture {shortUrl}
	UrlReWithID  = regexp.MustCompile(`^/([^/]+)$`)
	UrlRePrivate = regexp.MustCompile(`^/private/([^/]+)$`)
	UrlRePublic  = regexp.MustCompile(`^/public/([^/]+)$`)
)
