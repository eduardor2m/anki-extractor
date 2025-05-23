package handlers

import (
	"anki-extractor/internal/services"
	"anki-extractor/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type AnkiHandler struct {
	ankiService services.AnkiService
}

func NewAnkiHandler(s services.AnkiService) *AnkiHandler {
	return &AnkiHandler{ankiService: s}
}

func (ah *AnkiHandler) UploadAnkiCollectionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		return
	}

	file, header, err := r.FormFile("anki_file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter arquivo: %v", err), http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Erro ao fechar arquivo temporário: %v", err)
		} else {
			log.Printf("Arquivo temporário fechado com sucesso.")
		}
	}(file)

	uploadedFileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	uploadedFilePath := filepath.Join("tmp", uploadedFileName)

	outFile, err := os.Create(uploadedFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar arquivo temporário de upload: %v", err), http.StatusInternalServerError)
		return
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			log.Printf("Erro ao fechar arquivo temporário %s: %v", outFile.Name(), err)
		} else {
			log.Printf("Arquivo temporário %s fechado com sucesso.", outFile.Name())
		}
	}(outFile)

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao salvar arquivo de upload: %v", err), http.StatusInternalServerError)
		return
	}

	var ankiDBPath string
	if filepath.Ext(uploadedFilePath) == ".apkg" {

		unzipDest := filepath.Join("tmp", fmt.Sprintf("unzipped-%d", time.Now().UnixNano()))
		if err := os.MkdirAll(unzipDest, 0755); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao criar diretório para descompactação: %v", err), http.StatusInternalServerError)
			return
		}

		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				log.Printf("Erro ao remover diretório temporário %s: %v", path, err)
			} else {
				log.Printf("Diretório temporário %s removido com sucesso.", path)
			}
		}(unzipDest)

		log.Printf("Descompactando .apkg de %s para %s", uploadedFilePath, unzipDest)
		ankiDBPath, err = utils.Unzip(uploadedFilePath, unzipDest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao descompactar .apkg: %v", err), http.StatusInternalServerError)
			return
		}
		log.Printf("Arquivo .anki21/.anki2 encontrado em: %s", ankiDBPath)

	} else if filepath.Ext(uploadedFilePath) == ".anki21" || filepath.Ext(uploadedFilePath) == ".anki2" {

		ankiDBPath = uploadedFilePath
	} else {

		http.Error(w, "Tipo de arquivo não suportado. Por favor, envie um .anki21, .anki2 ou .apkg.", http.StatusBadRequest)
		return
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("Erro ao remover arquivo temporário %s: %v", name, err)
		} else {
			log.Printf("Arquivo temporário %s removido com sucesso.", name)
		}
	}(uploadedFilePath)

	flashcards, err := ah.ankiService.GetFlashcardsWithDeckNames(ankiDBPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar arquivo Anki: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flashcards)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar resposta JSON: %v", err), http.StatusInternalServerError)
		return
	}
}
