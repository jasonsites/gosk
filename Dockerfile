# syntax=docker/dockerfile:1

FROM golang:1.21-alpine3.18 AS build
WORKDIR /src
COPY go.* package.json ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o bin/server cmd/httpserver/main.go


FROM alpine:3.18
WORKDIR /app
COPY --from=build /src/package.json /src/config/config.toml /src/bin/server /app/
COPY --from=build /src/database/migrations /app/database/migrations

ENV MIGRATE_VERSION=v4.15.2
RUN apk --no-cache add curl \
  # migrate
  && curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz \
  && mv migrate /usr/bin/migrate \
  && chmod +x /usr/bin/migrate \
  && rm LICENSE README.md

EXPOSE 9202
CMD [ "/app/server" ]
