package dockercontainer

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

type containerBuilder interface {
	// setFields устанавливает все поля билдера, которые будут являться переменными окружения при создании контейнера
	setFields() error

	// setBase устанавливает контекст и докер клиент для билдера
	setBase(ctx context.Context, cli *client.Client)

	// getBase возвращает докер клиента и контекст билдера
	getBase() (*client.Client, context.Context)

	// makeContainer создает контейнер
	makeContainer() (container.ContainerCreateCreatedBody, error)
}

// containerDirector структура директора
type containerDirector struct {
	builder containerBuilder // поле билдера
}

// containerBase базовая структура, поля которой должны быть у каждого билдера
type containerBase struct {
	cli *client.Client  // поле докер клиента
	ctx context.Context // контекст
}

// constructContainer конструирует контейнер конкретного билдера
func (cd *containerDirector) constructContainer() error {
	if err := cd.builder.setFields(); err != nil {
		return err
	}

	cli, ctx := cd.builder.getBase()

	resp, err := cd.builder.makeContainer()
	if err != nil {
		if errdefs.IsConflict(errors.Unwrap(err)) {
			// TODO: если контейнер существует и он выключен, то нужно его стартануть, а не просто выйти из функции
			return nil
		}
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

// newContainerDirector создает нового директора
func newContainerDirector(b containerBuilder) (*containerDirector, error) {
	return &containerDirector{
		builder: b,
	}, nil
}

// setContainerBuilder устанавливает билдера для директора
func (cd *containerDirector) setContainerBuilder(b containerBuilder) {
	cd.builder = b
}

// TODO: контейнер для мигрирования в базу

// StartAllContainers создает и запускает все контейнеры
func StartAllContainers() error {
	vb := &vaultBuilder{}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	ctx := context.TODO()

	vb.setBase(ctx, cli)

	d, err := newContainerDirector(vb)
	if err != nil {
		return err
	}

	if err := d.constructContainer(); err != nil {
		return err
	}

	pb := &postgresBuilder{}

	pb.setBase(ctx, cli)
	d.setContainerBuilder(pb)

	if err := d.constructContainer(); err != nil {
		return err
	}

	return nil
}
