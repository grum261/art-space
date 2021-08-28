package main

import (
	"art_space/internal/rest"
	"encoding/json"
	"flag"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

func main() {
	var output string

	flag.StringVar(&output, "path", "", "Путь дла генерации OpenAPI 3 файлов")
	flag.Parse()

	if output == "" {
		logrus.Fatalf("Не передан путь")
	}

	swagger := rest.NewOpenApi3()

	data, err := json.Marshal(&swagger)
	if err != nil {
		logrus.Fatalf("Не удалось замаршалить json: %v", err)
	}

	if err := os.WriteFile(path.Join(output, "openapi3.json"), data, 0664); err != nil {
		logrus.Fatalf("Не удалось записать json в файл: %s", err)
	}

	data, err = yaml.Marshal(&swagger)
	if err != nil {
		logrus.Fatalf("Не удалось замаршалить yaml: %s", err)
	}

	if err := os.WriteFile(path.Join(output, "openapi3.yaml"), data, 0664); err != nil {
		logrus.Fatalf("Не удалось записать yaml в файл: %s", err)
	}

	logrus.Println("Все сгенерировано")
}
