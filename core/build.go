package core

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// tarProjectContext cria um tarball da pasta do projeto (Dockerfile incluso)
func tarProjectContext(projectPath string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err := filepath.Walk(projectPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(projectPath, file)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	tw.Close()
	return buf, nil
}

type buildMessage struct {
	Stream string          `json:"stream,omitempty"`
	Error  string          `json:"error,omitempty"`
	Aux    json.RawMessage `json:"aux,omitempty"`
}

type auxPayload struct {
	ID string `json:"ID"`
}

func extractImageId(buildOutput io.Reader) (string, error) {
	decoder := json.NewDecoder(buildOutput)

	for decoder.More() {
		var msg buildMessage

		if err := decoder.Decode(&msg); err != nil {
			return "", fmt.Errorf("erro ao decodificar saída do build: %w", err)
		}

		if msg.Error != "" {
			return "", fmt.Errorf("erro no build: %s", msg.Error)
		}

		if len(msg.Aux) > 0 {
			var payload auxPayload
			if err := json.Unmarshal(msg.Aux, &payload); err != nil {
				return "", fmt.Errorf("erro ao extrair ID da imagem: %w", err)
			}

			return payload.ID, nil
		}
	}

	return "", fmt.Errorf("erro ao extrair ID da imagem")
}