package vault

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

type Provider struct {
	path    string
	client  *api.Logical
	results map[string]map[string]string
}

func New(token, addr, path string) (*Provider, error) {
	client, err := api.NewClient(
		&api.Config{
			Address: addr,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к клиенту хранилища: %w", err)
	}

	client.SetToken(token)

	return &Provider{
		path:    path,
		client:  client.Logical(),
		results: make(map[string]map[string]string),
	}, nil
}

func (p *Provider) WriteSecret(path string, data map[string]interface{}) error {
	if _, err := p.client.Write(p.path+"/data/"+path, data); err != nil {
		return err
	}

	return nil
}

func (p *Provider) Get(v string) (string, error) {
	split := strings.Split(v, ":")
	if len(split) == 1 {
		return "", fmt.Errorf("отсутствует значение ключа")
	}

	pathSecret, key := split[0], split[1]

	res, ok := p.results[pathSecret]
	if ok {
		val, ok := res[key]
		if !ok {
			return "", errors.New("ключ не найден в закэшированных данных")
		}

		return val, nil
	}

	secret, err := p.client.Read(fmt.Sprintf("%s/data/%s", p.path, pathSecret))
	if err != nil {
		return "", fmt.Errorf("ошибка чтения из хранилища: %w", err)
	}

	if secret == nil {
		return "", errors.New("секрет не найден")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("некорректные данные в секрете")
	}

	secrets := make(map[string]string)

	for k, v := range data {
		secrets[k] = fmt.Sprint(v)
	}

	val, ok := secrets[key]
	if !ok {
		return "", errors.New("ключ не найден в полученных данных")
	}

	p.results[pathSecret] = secrets

	return val, nil
}
