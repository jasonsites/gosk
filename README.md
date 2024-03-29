# gosk
Go Starter Kit for Modular Backend Applications

## Documentation
- [Architecture](./documentation/architecture.md)
- [Getting Started](./documentation/getting-started.md)

## Installation
Clone the repository
```sh
$ git clone git@github.com:jasonsites/gosk.git
$ cd gosk
```

## Development
**Prerequisites**
- *[Docker Desktop](https://www.docker.com/products/docker-desktop)*
- *[Go 1.21+](https://golang.org/doc/install)*

**Show all commands**
```sh
$ docker compose run --rm api just
```

### Migrations
**Run all up migrations**
```sh
$ docker compose run --rm api just migrate
```

**Run up migrations {n} steps**
```sh
$ docker compose run --rm api just migrate-up svcdb {n}
```

**Run down migrations {n} steps**
```sh
$ docker compose run --rm api just migrate-down svcdb {n}
```

**Create new migration**
```sh
$ docker compose run --rm api just migrate-create {name}
```

### Server
**Run http server in development mode**
```sh
$ docker compose run --rm --service-ports api
```

### Testing
**Run unit tests**
```sh
$ docker compose run --rm api just test-unit
```

**Run integration tests**
```sh
$ docker compose run --rm api just test-int
```

**Run full test suite with code coverage**
```sh
$ docker compose run --rm api just coverage
```

## Building
**Compile server binary**
```sh
$ go build -mod vendor -o out/bin/domain ./cmd/httpserver
```

## License
Copyright (c) 2022 Jason Sites

Licensed under the [MIT License](LICENSE.md)
