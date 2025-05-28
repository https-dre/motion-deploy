package core

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"motion/pkgs/models"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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

func (c *Instance) buildImage(ctx context.Context, projectPath, imageName string) (string, error) {
	cli := c.Client
	// Empacota a pasta do projeto como contexto
	tarContext, err := tarProjectContext(projectPath)
	fmt.Printf("Packing %s\n", imageName)

	if err != nil {
		return "", fmt.Errorf("erro ao criar tar do contexto: %w", err)
	}

	// builda a imagem
	fmt.Println("Creating Image...")
	buildResp, err := cli.ImageBuild(ctx, tarContext, types.ImageBuildOptions{
		Tags:       []string{imageName},
		Remove:     true,
		Dockerfile: "Dockerfile", // padrão
	})

	if err != nil {
		return "", fmt.Errorf("erro ao buildar imagem: %w", err)
	}

	defer buildResp.Body.Close()

	imageID, err := extractImageId(buildResp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("Image created [ID = %s]\n", imageID)

	return imageID, nil
}

func (c *Instance) createContainer(ctx context.Context, service *models.CoreApplication) error {
	cli := c.Client

	imageList, err := cli.ImageList(ctx, image.ListOptions{})

	if err != nil {
		return err
	}

	for _, image := range imageList {
		if image.ID == service.ImageId {
			containerResponse, err := cli.ContainerCreate(ctx, &container.Config{
				Image: service.ImageId,
				Env:   service.Env,
			}, nil, nil, nil, service.Name)

			if err != nil {
				return err
			}

			service.CoreId = containerResponse.ID

			return nil
		}
	}

	return nil
}

func (c *Instance) startContainer(ctx context.Context, containerID string) error {
	cli := c.Client

	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container with id %s, %v", containerID, err)
	}

	return nil
}