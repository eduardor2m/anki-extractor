package routes

import (
	"anki-extractor/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterAnkiRoutes(r *mux.Router, ah *handlers.AnkiHandler) {
	r.HandleFunc("/upload-anki", ah.UploadAnkiCollectionHandler).Methods("POST")
}
