set shell := ["bash", "-uc"]

# Defaults ========================================================================================
project := 'domain'

# Commands ========================================================================================
# show this help
help:
  just --list

# remove build related files
clean:
  rm -rf bin
  rm -rf out
  rm -f test/coverage/profile.cov

# Migrations ======================================================================================
# migrate down
migrate-down db +step='-all':
  migrate -path ./database/migrations -database postgres://postgres:postgres@postgres_db:5432/{{db}}?sslmode=disable down {{step}}

# migrate up
migrate-up db *step:
  migrate -path ./database/migrations -database postgres://postgres:postgres@postgres_db:5432/{{db}}?sslmode=disable up {{step}}

# migrate up -all (alias)
migrate:
  just migrate-up svcdb

# create migration with {{name}}
migrate-create name:
	migrate create -ext sql -dir ./database/migrations -format unix {{name}}

# Run =============================================================================================
# run {http}server in dev mode
serve-dev +protocol='http':
  go run ./cmd/{{protocol}}server/main.go

# run {http}server in dev mode with file monitor
serve +protocol='http':
  air --tmp_dir="out/tmp" --build.cmd="go build -mod readonly -o out/tmp/domain ./cmd/{{protocol}}server" --build.bin="out/tmp/domain"

# Test ============================================================================================
# run integration tests
test-int:
  go test -v ./test/inegration/...

# run unit tests
test-unit:
  go test -v ./internal/...

# run unit and inegration tests with coverage report
coverage:
  just migrate-up testdb
  go build -cover -o ./test/bin/domain ./cmd/httpserver
  ./test/bin/domain
  go test -cover ./... -args -test.gocoverdir="$PWD/test/coverage/unit"

covprint:
  go tool covdata percent -i=./test/coverage/integration,./test/coverage/unit

covprofile:
  go tool covdata textfmt -i=./test/coverage/integration,./test/coverage/unit -o test/coverage/profile

# html coverage report
covreport:
  go tool cover -html=./coverage/profile.cov
