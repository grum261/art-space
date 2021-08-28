package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	fmt.Println(migrate.ErrDirty{})
}
