package dockercontainer

import (
	"art_space/internal/envvar"
	"context"
	"fmt"
	"net/url"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// vaultBuilder конкретный тип билдера для контейнера vault
type vaultBuilder struct {
	containerBase // базовая структура, содержащая докер клиент и контекст
}

// getBase возвращает докер клиента и контекст билдера
func (vb *vaultBuilder) getBase() (*client.Client, context.Context) {
	return vb.cli, vb.ctx
}

// setBase устанавливает контекст и докер клиент для билдера
func (vb *vaultBuilder) setBase(ctx context.Context, cli *client.Client) {
	vb.cli = cli
	vb.ctx = ctx
}

func (vb *vaultBuilder) makeContainer() (container.ContainerCreateCreatedBody, error) {
	u, err := url.ParseRequestURI(envvar.Config.Vault.Address)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("(dockercontainer.makeContainer) ошибка парсинга адреса Vault: %w", err)
	}

	resp, err := vb.cli.ContainerCreate(
		vb.ctx, &container.Config{
			ExposedPorts: nat.PortSet{"8300": struct{}{}},
			Env: []string{
				fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%s", envvar.Config.Vault.Token),
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
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("(dockercontainer.makeContainer) ошибка при создании контейнера Vault: %w", err)
	}

	return resp, nil
}
