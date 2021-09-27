FROM golang:1.17.1-alpine3.14 as builder

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w" \
    art_space/cmd/rest-server

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -ldflags "-s -w" -tags 'pgx' \
      github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.0

#-

FROM alpine:3.14 AS certificates

RUN apk --no-cache add ca-certificates

#-

FROM scratch

WORKDIR /api/

ENV PATH=/api/bin/:$PATH

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/rest-server ./bin/rest-server
COPY --from=builder /go/bin/migrate ./bin/migrate
COPY --from=builder /build/.env.dev .
COPY --from=builder /build/internal/pgdb/sql .

EXPOSE 8000

CMD ["rest-server", "-env", "/api/.env.dev"]