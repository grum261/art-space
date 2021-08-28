package dockercontainer

import (
	"context"

	"github.com/docker/docker/client"
)

// containerBase базовая структура, поля которой должны быть у каждого билдера
type containerBase struct {
	cli *client.Client  // поле докер клиента
	ctx context.Context // контекст
}
