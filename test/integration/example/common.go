package exampletest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk/internal/modules/example"
	"github.com/jasonsites/gosk/internal/resolver"
	utils "github.com/jasonsites/gosk/test/testutils"
)

type Suite struct {
	DB          *pgxpool.Pool
	Handler     http.Handler
	Method      string
	Resolver    *resolver.Resolver
	RoutePrefix string
}

func (s *Suite) SetupSuite(tb testing.TB) func(tb testing.TB) {
	conf := &resolver.Config{}
	resolver, err := utils.InitializeResolver(conf, "http")
	if err != nil {
		tb.Fatalf("app initialization error: %+v\n", err)
	}

	s.DB = resolver.PostgreSQLClient()
	s.Handler = resolver.HTTPServer().Server.Handler
	s.Method = http.MethodPost
	s.Resolver = resolver
	s.RoutePrefix = "/domain/examples"

	return func(tb testing.TB) {
		// teardown for test table
	}
}

func (s *Suite) SetupTest(tb testing.TB) func(tb testing.TB) {
	// setup for each test

	return func(tb testing.TB) {
		utils.Cleanup(s.Resolver)
	}
}

// insertExampleRecord inserts a db record for use in test setup
func insertExampleRecord(e *example.ExampleEntity, db *pgxpool.Pool) (*example.ExampleEntity, error) {
	var (
		statement    = "INSERT INTO %s %s VALUES %s RETURNING id"
		name         = "example_entity"
		insertFields = "(title,description,created_context)"
		values       = "($1,$2,$3)"
		query        = fmt.Sprintf(statement, name, insertFields, values)
	)

	var (
		title          = e.Title
		description    = e.Description.String
		createdContext = e.CreatedContext
	)

	// create new entity for db row scan and execute query
	entity := &example.ExampleEntity{}
	if err := db.QueryRow(
		context.Background(),
		query,
		title,
		description,
		createdContext,
	).Scan(
		&entity.ID,
	); err != nil {
		return nil, err
	}

	return entity, nil
}
