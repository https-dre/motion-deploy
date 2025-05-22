package core

import "github.com/docker/docker/client"

type DockerClient = *client.Client

type Instance struct {
	Client       DockerClient
	Applications []CoreApplication
}

type CoreApplication struct {
	CoreId  string
	Name    string
	Env     []string
	ImageId string
}

func NewDockerClient() DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return cli
}

func NewCore() *Instance {
	cli := NewDockerClient()

	return &Instance{
		Client: cli,
	}
}

func (c *Instance) ListApplications() []CoreApplication {
	return []CoreApplication{}
}
