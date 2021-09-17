package envvar

import (
	"art_space/internal/envvar/vault"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

//go:generate counterfeiter -o envvartesting/provider.gen.go . Provider

var (
	Config        = &Configuration{}
	VaultProvider = &Vault{}
)

type Configuration struct {
	DB    dBConfiguration
	Vault vaultConfiguration
}

type vaultConfiguration struct {
	Path    string
	Address string
	Token   string
}

type dBConfiguration struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
}

type Provider interface {
	Get(key string) (string, error)
}

type Vault struct {
	provider Provider
}

func Load(fileName string) error {
	if err := godotenv.Load(fileName); err != nil {
		return fmt.Errorf("ошибка при загрузке env файла: %w", err)
	}

	return nil
}

func New(prodiver Provider) *Vault {
	return &Vault{
		provider: prodiver,
	}
}

func (c *Vault) Get(key string) (string, error) {
	res := os.Getenv(key)
	valSecret := os.Getenv(fmt.Sprintf("%s_SECURE", key))

	if valSecret != "" {
		valSecretRes, err := c.provider.Get(valSecret)
		if err != nil {
			return "", err
		}

		res = valSecretRes
	}

	return res, nil
}

func newVaultConfig() error {
	var env string

	flag.StringVar(&env, "env", "", "Название файла с переменными окружения")
	flag.Parse()

	if err := Load(env); err != nil {
		return fmt.Errorf("(envvar.Load) не удалось загрузить конфигурацию: %w", err)
	}

	Config.Vault.Path = os.Getenv("VAULT_PATH")
	Config.Vault.Token = os.Getenv("VAULT_TOKEN")
	Config.Vault.Address = os.Getenv("VAULT_ADDRESS")

	p, err := vault.New(Config.Vault.Token, Config.Vault.Address, Config.Vault.Path)
	if err != nil {
		return fmt.Errorf("(envvar.vault.New) не удалось создать провайдер: %w", err)
	}

	VaultProvider = New(p)

	return nil
}

func newDBConfig() {
	get := func(v string) string {
		res, err := VaultProvider.Get(v)
		if err != nil {
			log.Fatalf("(envvar.Get) не удалось получить значение конфигурации для %s: %v", v, err)
		}

		return res
	}

	Config.DB.Username = get("PGDB_USERNAME")
	Config.DB.Password = get("PGDB_PASSWORD")
	Config.DB.Host = get("PGDB_HOST")
	Config.DB.Port = get("PGDB_PORT")
	Config.DB.Name = get("PGDB_NAME")
}

func init() {
	if err := newVaultConfig(); err != nil {
		log.Fatal(err)
	}

	newDBConfig()
}
