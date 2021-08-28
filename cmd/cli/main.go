package main

import (
	"art_space/pkg/openapi3"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	client, err := openapi3.NewClientWithResponses("http://0.0.0.0:9234")
	if err != nil {
		logrus.Fatalf("Couldn't instantiate client: %s", err)
	}

	fmt.Println(client)
}
