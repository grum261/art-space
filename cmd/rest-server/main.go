package main

import (
	dockercontainer "art_space/internal/docker-container"
	"art_space/internal/envvar"
	"art_space/internal/models/service"
	"art_space/internal/pgdb"
	"art_space/internal/rest"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.TODO()

	// TODO: сделать флаг для запуска без докера
	if err := dockercontainer.StartAllContainers(); err != nil {
		logrus.Fatal(err)
	}

	// TODO: добавить мигрирование

	db := pgdb.NewDB(ctx, envvar.Configuration)
	defer db.Close(ctx)

	postRepo := pgdb.NewPost(db)
	svc := service.NewPost(postRepo)

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
		IdleTimeout:  time.Second * 1,
	})

	rest.RegisterOpenApi(app)
	rest.NewPostHandler(svc).RegisterRoutes(app)

	app.Static("/", "./assets/swagger-ui")

	logrus.Fatal(app.Listen(":8000"))
}
