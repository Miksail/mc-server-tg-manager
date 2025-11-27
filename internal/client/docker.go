package client

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"mc-server-tg-manager/internal/model"
)

type DockerClient struct {
	cli           *client.Client
	containerName string
}

func NewDockerClient(containerName string) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerClient{
		cli:           cli,
		containerName: containerName,
	}, nil
}

func (d *DockerClient) Status() (model.ServerStatus, error) {
	insp, err := d.cli.ContainerInspect(context.Background(), d.containerName)
	if err != nil {
		return model.ServerStatusUnknown, err
	}
	if !insp.State.Running {
		return model.ServerStatusStopped, nil
	}
	return model.ServerStatusRunning, nil
}

func (d *DockerClient) Start() error {
	return d.cli.ContainerStart(context.Background(), d.containerName, container.StartOptions{})
}
