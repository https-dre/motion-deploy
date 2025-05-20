package core

import "github.com/docker/docker/client"


type DockerClient = *client.Client

type CoreInstance struct {
	Client DockerClient
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return cli
}

func NewCore() *CoreInstance {
	cli := NewDockerClient()

	return &CoreInstance{
		Client: cli,
	}
}