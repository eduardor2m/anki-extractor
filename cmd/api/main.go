package main

import (
	"anki-extractor/config"
	"anki-extractor/internal/handlers"
	repositories "anki-extractor/internal/repository"
	"anki-extractor/internal/routes"
	"anki-extractor/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		err := os.Mkdir("tmp", 0755)
		if err != nil {
			log.Fatalf("Falha ao criar diretório temporário 'tmp/': %v", err)
		}
	}

	cfg := config.LoadConfig()

	ankiRepo := repositories.NewAnkiRepository()
	ankiService := services.NewAnkiService(ankiRepo)
	ankiHandler := handlers.NewAnkiHandler(ankiService)

	r := mux.NewRouter()

	routes.RegisterAnkiRoutes(r, ankiHandler)

	port := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("API Anki rodando em http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
