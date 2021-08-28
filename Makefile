include .env.dev

run:
	go run cmd/rest-server/main.go -env .env.dev

build:
	go build -o art_space cmd/main.go

migrate-up:
	migrate -path=./migrations -database=postgres://${PGDB_USERNAME}:${PGDB_PASSWORD}@${PGDB_HOST}:${PGDB_PORT}/${PGDB_NAME}?sslmode=disable up

migrate-down:
	migrate -path=migrations -database=postgres://${PGDB_USERNAME}:${PGDB_PASSWORD}@${PGDB_HOST}:${PGDB_PORT}/${PGDB_NAME} down

gen-sql:
	sqlc generate -f "internal/pgdb/sqlc.yml"