package lib_testcontainers

import (
	"maps"

	"github.com/docker/go-connections/nat"
)

type (
	containerConfig struct {
		Name         string
		Image        string
		Env          map[string]string
		ExposedPorts []string
		WaitPort     nat.Port
	}

	// Option -
	Option func(*containerConfig)
)

// WithName -
func WithName(name string) Option {
	return func(cfg *containerConfig) {
		cfg.Name = name
	}
}

// WithImage -
func WithImage(image string) Option {
	return func(cfg *containerConfig) {
		cfg.Image = image
	}
}

// WithEnv -
func WithEnv(env map[string]string) Option {
	return func(cfg *containerConfig) {
		maps.Copy(cfg.Env, env)
	}
}

// WithExposedPort -
func WithExposedPort(port string) Option {
	return func(cfg *containerConfig) {
		cfg.ExposedPorts = append(cfg.ExposedPorts, port)
		p, err := nat.NewPort("tcp", port)
		_ = err
		cfg.WaitPort = p
	}
}
