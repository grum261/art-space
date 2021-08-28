package dockercontainer

import (
	"art_space/internal/envvar"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// postgresBuilder конкретный тип билдера для контейнера postgres
type postgresBuilder struct {
	containerBase                     // базовая структура, содержащая докер клиент и контекст
	username, password, dbName string // поля переменных окружения контейнера
}

// getBase возвращает докер клиента и контекст билдера
func (pb *postgresBuilder) getBase() (*client.Client, context.Context) {
	return pb.cli, pb.ctx
}

func (pb *postgresBuilder) setFields() (err error) {
	get := func(v string) (string, error) {
		res, err := envvar.Configuration.Get(v)
		if err != nil {
			return "", fmt.Errorf("не удалось получить значение конфигурации для %s: %w", v, err)
		}

		return res, nil
	}

	pb.username, err = get("PGDB_USERNAME")
	if err != nil {
		return err
	}

	pb.password, err = get("PGDB_PASSWORD")
	if err != nil {
		return err
	}

	pb.dbName, err = get("PGDB_NAME")
	if err != nil {
		return err
	}

	return nil
}

// setBase устанавливает контекст и докер клиент для билдера
func (pb *postgresBuilder) setBase(ctx context.Context, cli *client.Client) {
	pb.cli = cli
	pb.ctx = ctx
}

func (pb *postgresBuilder) makeContainer() (container.ContainerCreateCreatedBody, error) {
	resp, err := pb.cli.ContainerCreate(
		pb.ctx,
		&container.Config{
			ExposedPorts: nat.PortSet{"5432": struct{}{}},
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", pb.username),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", pb.password),
				fmt.Sprintf("POSTGRES_DB=%s", pb.dbName),
			},
			Image: "postgres:latest",
		}, &container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				nat.Port("5432"): {{HostIP: "0.0.0.0", HostPort: "5432"}},
			},
			RestartPolicy: container.RestartPolicy{Name: "always"},
		}, nil, nil, "art_space_postgres",
	)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("(makeContainer) не удалось создать контейнер Postgres: %w", err)
	}

	return resp, nil
}
