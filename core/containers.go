package core

import (
	"context"
	"fmt"
	"motion/pkgs/config"
	"motion/pkgs/models"
	"strings"

	"github.com/docker/docker/api/types/container"
)

// builda a imagem local e sobe um container
func (c *Instance) BuildAndRun(projectPath, imageName, containerName string) (models.CoreApplication, error) {
	ctx := context.Background()

	imageId, err := c.buildImage(ctx, projectPath, imageName)

	if err != nil {
		return models.CoreApplication{}, err
	}

	if imageId == "" {
		return models.CoreApplication{}, fmt.Errorf("ImageID is invalid")
	}

	newServiceApp := models.CoreApplication{
		Name:    containerName,
		ImageId: imageId,
	}

	fmt.Printf("Creating container %s\n", newServiceApp.Name)
	if err = c.createContainer(ctx, &newServiceApp); err != nil {
		return models.CoreApplication{}, fmt.Errorf("an error ocurred in container creation: %v", err)
	}

	fmt.Printf("Starting container %s with Id %s\n", newServiceApp.Name, newServiceApp.CoreId)
	if err := c.startContainer(ctx, newServiceApp.CoreId); err != nil {
		return models.CoreApplication{}, err
	}

	return newServiceApp, nil
}

func (c *Instance) ListContainers() []models.CoreApplication {
	var ctx context.Context = context.Background()
	containers_list, err := c.Client.ContainerList(ctx, container.ListOptions{})

	if err != nil {
		fmt.Println("An error ocurred while listing containers_list")
		fmt.Println(err)
	}

	var motion_services []models.CoreApplication

	for _, container := range containers_list {
		for _, repo := range config.Repos {
			if strings.Contains(container.ID, repo.Service.CoreId) {
				motion_services = append(motion_services, repo.Service)
			}
		}
	}

	return motion_services
}
