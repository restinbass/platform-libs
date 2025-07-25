package lib_testcontainers

import (
	"context"
	"errors"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Container -
type Container struct {
	Container  testcontainers.Container
	host       string
	mappedPort string
}

// RunContainer -
func RunContainer(ctx context.Context, opts ...Option) (*Container, error) {
	config := &containerConfig{
		Env: make(map[string]string),
	}
	for _, opt := range opts {
		opt(config)
	}

	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Name:         config.Name,
		Env:          config.Env,
		ExposedPorts: config.ExposedPorts,
		WaitingFor:   wait.ForListeningPort(config.WaitPort),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := c.Host(ctx)
	if err != nil {
		terminateErr := c.Terminate(ctx)
		return nil, errors.Join(err, terminateErr)
	}
	mappedPort, err := c.MappedPort(ctx, config.WaitPort)
	if err != nil {
		terminateErr := c.Terminate(ctx)
		return nil, errors.Join(err, terminateErr)
	}

	return &Container{
		Container:  c,
		host:       host,
		mappedPort: mappedPort.Port(),
	}, nil
}

// Получить host
func (c *Container) Host() string {
	return c.host
}

// Получить порт
func (c *Container) Port() string {
	return c.mappedPort
}

// Остановить контейнер
func (c *Container) Terminate(ctx context.Context) error {
	return c.Container.Terminate(ctx)
}
