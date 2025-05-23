package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(src, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo zip: %w", err)
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			fmt.Printf("erro ao fechar arquivo zip: %v\n", err)
		} else {
			fmt.Println("arquivo zip fechado com sucesso.")
		}
	}(r)

	var foundAnki2Path string

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return "", fmt.Errorf("%s: caminho inválido de arquivo", fpath)
		}

		if f.FileInfo().IsDir() {

			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return "", err
			}
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", fmt.Errorf("erro ao criar arquivo %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			err := outFile.Close()
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf("erro ao abrir arquivo no zip %s: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return "", fmt.Errorf("erro ao copiar arquivo %s: %w", f.Name, err)
		}

		err = rc.Close()
		if err != nil {
			return "", err
		}
		err = outFile.Close()
		if err != nil {
			return "", err
		}

		if strings.HasSuffix(strings.ToLower(f.Name), ".anki21") {
			return fpath, nil
		}
		if strings.HasSuffix(strings.ToLower(f.Name), ".anki21b") {
			return fpath, nil
		}
		if strings.HasSuffix(strings.ToLower(f.Name), ".anki2") {
			foundAnki2Path = fpath
		}
	}

	if foundAnki2Path != "" {
		return foundAnki2Path, nil
	}

	return "", fmt.Errorf("arquivo .anki21 ou .anki2 não encontrado dentro do pacote Anki")
}
