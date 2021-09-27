include .env.dev

run:
	go run cmd/rest-server/main.go -env .env.dev

build:
	go build -o art_space cmd/main.go

migrate-up:
	docker-compose run api migrate -path=./migrations -database=pgx://${PGDB_USERNAME}:${PGDB_PASSWORD}@localhost:${PGDB_PORT}/${PGDB_NAME}?sslmode=disable up
