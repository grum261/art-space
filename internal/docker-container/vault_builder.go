package dockercontainer

import (
	"art_space/internal/envvar"
	"art_space/internal/envvar/vault"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// vaultBuilder конкретный тип билдера для контейнера vault
type vaultBuilder struct {
	containerBase               // базовая структура, содержащая докер клиент и контекст
	token, path, address string // поля переменных окружения контейнера
}

// getBase возвращает докер клиента и контекст билдера
func (vb *vaultBuilder) getBase() (*client.Client, context.Context) {
	return vb.cli, vb.ctx
}

// setFields устанавливает все поля билдера, которые будут являться переменными окружения при создании контейнера
func (vb *vaultBuilder) setFields() error {
	var env string

	flag.StringVar(&env, "env", "", "Название файла с переменными окружения")
	flag.Parse()

	if err := envvar.Load(env); err != nil {
		return fmt.Errorf("не удалось загрузить конфигурацию: %w", err)
	}

	vb.path = os.Getenv("VAULT_PATH")
	vb.token = os.Getenv("VAULT_TOKEN")
	vb.address = os.Getenv("VAULT_ADDRESS")

	p, err := vault.New(vb.token, vb.address, vb.path)
	if err != nil {
		return fmt.Errorf("(setFields) не удалось создать провайдер: %w", err)
	}

	envvar.Configuration = envvar.New(p)

	return nil
}

// setBase устанавливает контекст и докер клиент для билдера
func (vb *vaultBuilder) setBase(ctx context.Context, cli *client.Client) {
	vb.cli = cli
	vb.ctx = ctx
}

func (vb *vaultBuilder) makeContainer() (container.ContainerCreateCreatedBody, error) {
	u, err := url.ParseRequestURI(vb.address)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("(makeContainer) ошибка парсинга адреса Vault: %w", err)
	}

	resp, err := vb.cli.ContainerCreate(
		vb.ctx, &container.Config{
			ExposedPorts: nat.PortSet{"8300": struct{}{}},
			Env: []string{
				fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%s", vb.token),
				fmt.Sprintf("VAULT_DEV_LISTEN_ADDRESS=%s", u.Host),
			},
			Image: "vault:latest",
		}, &container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				nat.Port("8300"): {{HostIP: "0.0.0.0", HostPort: "8300"}},
			},
			CapAdd:        []string{"IPC_LOCK"},
			RestartPolicy: container.RestartPolicy{Name: "always"},
		}, nil, nil, "art_space_vault",
	)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("(makeContainer) ошибка при создании контейнера Vault: %w", err)
	}

	return resp, nil
}
