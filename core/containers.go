package core

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

// builda a imagem local e sobe um container
func (c *Instance) BuildAndRunService(projectPath, imageName, containerName string) error {
	ctx := context.Background()

	imageId, err := c.buildImage(ctx, projectPath, imageName)

	if err != nil {
		return err
	}

	if imageId == "" {
		return fmt.Errorf("ImageID is invalid")
	}

	newServiceApp := CoreApplication{
		Name: containerName,
		ImageId: imageId,
	}

	err = c.createContainer(ctx, &newServiceApp)

	if err != nil {
		return fmt.Errorf("an error ocurred in container creation: %v", err)
	}

	if err := c.startContainer(ctx, newServiceApp.CoreId); err != nil {
		return err
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

func (c *Instance) createContainer(ctx context.Context, service *CoreApplication) error {
	cli := c.Client

	imageList, err := cli.ImageList(ctx, image.ListOptions{})

	if err != nil {
		return err
	}

	for _, image := range imageList {
		if image.ID == service.ImageId {
			containerResponse, err := cli.ContainerCreate(ctx, &container.Config{
				Image: image.RepoTags[0],
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

func (c *Instance) buildImage(ctx context.Context, projectPath, imageName string) (string, error) {
	cli := c.Client
	// Empacota a pasta do projeto como contexto
	tarContext, err := tarProjectContext(projectPath)
	if err != nil {
		return "", fmt.Errorf("erro ao criar tar do contexto: %w", err)
	}

	// builda a imagem
	buildResp, err := cli.ImageBuild(ctx, tarContext, types.ImageBuildOptions{
		Tags:       []string{imageName},
		Remove:     true,
		Dockerfile: "Dockerfile", // padrão
	})

	if err != nil {
		return "", fmt.Errorf("erro ao buildar imagem: %w", err)
	}

	defer buildResp.Body.Close()
	io.Copy(os.Stdout, buildResp.Body)

	var imageID string

	// Estrutura dos logs de build
	type buildLine struct {
		Stream string `json:"stream"`
	}

	decoder := json.NewDecoder(buildResp.Body)
	for decoder.More() {
		var msg buildLine
		if err := decoder.Decode(&msg); err != nil {
			return "", fmt.Errorf("erro ao decodificar saída do build: %w", err)
		}

		fmt.Print(msg.Stream) // mostra os logs no terminal

		if strings.HasPrefix(msg.Stream, "Successfully built") {
			parts := strings.Fields(msg.Stream)
			if len(parts) == 3 {
				imageID = strings.TrimSpace(parts[2])
			}
		}
	}

	if imageID == "" {
		return "", fmt.Errorf("ID da imagem não encontrado nos logs")
	}

	return "sha256:" + imageID, nil
}
