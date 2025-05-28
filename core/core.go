package core

import (
	"fmt"
	"motion/pkgs/models"

	"github.com/docker/docker/client"
)

type DockerClient = *client.Client

var Engine CoreInstance

type Instance struct {
	Client DockerClient
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("An error ocurred while create the docker integration")
		panic(err)
	}
	return cli
}

func NewDockerCore() CoreInstance {
	cli := NewDockerClient()

	return &Instance{
		Client: cli,
	}
}

func (c *Instance) ListApplications() []models.CoreApplication {
	return []models.CoreApplication{}
}
