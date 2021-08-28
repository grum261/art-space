package envvar

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

//go:generate counterfeiter -o envvartesting/provider.gen.go . Provider

var Configuration *VaultConfiguration

type Provider interface {
	Get(key string) (string, error)
}

type VaultConfiguration struct {
	provider Provider
}

func Load(fileName string) error {
	if err := godotenv.Load(fileName); err != nil {
		return fmt.Errorf("ошибка при загрузке env файла: %w", err)
	}

	return nil
}

func New(prodiver Provider) *VaultConfiguration {
	return &VaultConfiguration{
		provider: prodiver,
	}
}

func (c *VaultConfiguration) Get(key string) (string, error) {
	res := os.Getenv(key)
	valSecret := os.Getenv(fmt.Sprintf("%s_SECURE", key))

	if valSecret != "" {
		valSecretRes, err := c.provider.Get(valSecret)
		if err != nil {
			return "", fmt.Errorf("получение провайдером: %w", err)
		}

		res = valSecretRes
	}

	return res, nil
}
