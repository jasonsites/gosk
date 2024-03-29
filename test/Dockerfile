# syntax=docker/dockerfile:1

FROM golang:1.21
ENV MIGRATE_VERSION=v4.15.2
WORKDIR /app
COPY go.* package.json ./

RUN \
  # just
  curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin \
  && chmod +x /usr/local/bin/just \
  # air
  && curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
  && chmod +x install.sh \
  && sh install.sh \
  && mv ./bin/air /bin/air \
  && rm -rf ./bin install.sh \
  && mkdir -p out/tmp \
  # gotestsum
  && go install gotest.tools/gotestsum@latest \
  # migrate
  && curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz \
  && mv migrate /usr/bin/migrate \
  && chmod +x /usr/bin/migrate \
  && rm LICENSE README.md

COPY . .
RUN go mod download

EXPOSE 9202

CMD ["just"]
