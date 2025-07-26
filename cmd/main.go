package main

import (
	"context"
	"log"
	"time"

	lib_testcontainers "github.com/restinbass/platform-libs/pkg/testcontainers"
)

func main() {
	ctx := context.Background()

	postgresContainer, err := lib_testcontainers.RunContainer(
		ctx,
		lib_testcontainers.WithName("restinbass_payment_service_test_postgres"),
		lib_testcontainers.WithImage("postgres:16.4"),
		lib_testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_db",
		}),
		lib_testcontainers.WithExposedPort("5432"),
	)
	if err != nil {
		panic(err)
	}

	log.Printf("Host: %s, Port: %s\n", postgresContainer.Host(), postgresContainer.Port())
	time.Sleep(30 * time.Second)

	// Tests...
	if err := postgresContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	log.Printf("container terminated")
}
