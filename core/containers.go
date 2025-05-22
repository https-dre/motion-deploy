package core

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

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

	fmt.Printf("Creating container %s\n", newServiceApp.Name)
	if err = c.createContainer(ctx, &newServiceApp); err != nil {
		return fmt.Errorf("an error ocurred in container creation: %v", err)
	}

	fmt.Printf("Starting container %s with Id %s\n", newServiceApp.Name, newServiceApp.CoreId)
	if err := c.startContainer(ctx, newServiceApp.CoreId); err != nil {
		return err
	}

	return nil
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

func (c *Instance) createContainer(ctx context.Context, service *CoreApplication) error {
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


