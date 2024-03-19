package main

import (
	"net/http"

	"com.github/prajwalbharadwajbm/GoLinkify/internal/handlers"
	"com.github/prajwalbharadwajbm/GoLinkify/internal/storage"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &handlers.URLHandler{})
	store := storage.NewMemoryStorage()
	handlers.SetStorage(store)
	http.ListenAndServe(":8080", mux)
}
