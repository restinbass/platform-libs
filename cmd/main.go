package main

import (
	"context"
	"log"
	"time"

	lib_testcontainers "github.com/restinbass/platform-libs/pkg/testcontainers"
)

func main() {
	ctx := context.Background()
	env := map[string]string{
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "pass",
		"POSTGRES_DB":       "testdb",
	}

	c, err := lib_testcontainers.RunContainer(
		ctx,
		lib_testcontainers.WithName("testpg"),
		lib_testcontainers.WithImage("postgres:16"),
		lib_testcontainers.WithEnv(env),
		lib_testcontainers.WithExposedPort("5432"),
	)
	if err != nil {
		panic(err)
	}

	log.Printf("Host: %s, Port: %s\n", c.Host(), c.Port())
	time.Sleep(30 * time.Second)

	// Tests...
	if err := c.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	log.Printf("container terminated")
}
