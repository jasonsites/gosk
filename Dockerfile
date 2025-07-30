# syntax=docker/dockerfile:1

FROM golang:1.24-alpine3.22 AS base
WORKDIR /src
COPY go.* ./
RUN go mod download && go mod verify


FROM alpine:3.22 AS migrate-tools
ENV MIGRATE_VERSION=v4.18.3
RUN apk --no-cache add curl \
  && curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz \
  && mv migrate /usr/bin/migrate \
  && chmod +x /usr/bin/migrate \
  && rm LICENSE README.md


FROM base AS dev
WORKDIR /app
RUN apk --no-cache add curl bash \
  # just
  && curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin \
  && chmod +x /usr/local/bin/just \
  # air
  && curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
  && chmod +x install.sh \
  && sh install.sh \
  && mv ./bin/air /bin/air \
  && rm -rf ./bin install.sh \
  && mkdir -p out/tmp \
  # gotestsum
  && go install gotest.tools/gotestsum@latest

COPY --from=migrate-tools /usr/bin/migrate /usr/bin/migrate
COPY . .

EXPOSE 9202
CMD ["just"]


FROM base AS build
COPY . .
RUN CGO_ENABLED=0 go build -o bin/server cmd/httpserver/main.go


FROM alpine:3.22 AS prod
WORKDIR /app
COPY --from=build /src/config/config.toml /src/bin/server /app/
COPY --from=build /src/database/migrations /app/database/migrations
COPY --from=migrate-tools /usr/bin/migrate /usr/bin/migrate

EXPOSE 9202
CMD [ "/app/server" ]
