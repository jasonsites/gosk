package resolver

import (
	"context"
	"testing"
)

var r = NewResolver(context.Background(), nil)

func TestDefaultNewResolver(t *testing.T) {
	// non-nil
	if r.appContext == nil {
		t.Errorf("r.appContext expected to be non-nil; actual: %+v", r.appContext)
	}

	// nil
	if r.config != nil {
		t.Errorf("resolver.config expected to be nil; actual: %+v", r.config)
	}
	if r.domain != nil {
		t.Errorf("resolver.domain expected to be nil; actual: %+v", r.domain)
	}
	if r.exampleRepo != nil {
		t.Errorf("resolver.exampleRepo expected to be nil; actual: %+v", r.exampleRepo)
	}
	if r.exampleService != nil {
		t.Errorf("resolver.exampleService expected to be nil; actual: %+v", r.exampleService)
	}
	if r.httpServer != nil {
		t.Errorf("resolver.httpServer expected to be nil; actual: %+v", r.httpServer)
	}
	if r.log != nil {
		t.Errorf("resolver.log expected to be nil; actual: %+v", r.log)
	}
	if r.metadata != nil {
		t.Errorf("resolver.metadata expected to be nil; actual: %+v", r.metadata)
	}
	if r.postgreSQLClient != nil {
		t.Errorf("resolver.postgreSQLClient expected to be nil; actual: %+v", r.postgreSQLClient)
	}
}

func TestConfigComponent(t *testing.T) {
	actual := r.Config()
	if actual != r.config {
		t.Errorf("resolver.config expected to be a singleton; actual: %+v", actual)
	}
}

func TestDomainComponent(t *testing.T) {
	actual := r.Domain()
	if actual != r.domain {
		t.Errorf("resolver.domain expected to be a singleton; actual: %+v", actual)
	}
}

func TestExampleRepositoryComponent(t *testing.T) {
	actual := r.ExampleRepository()
	if actual != r.exampleRepo {
		t.Errorf("resolver.exampleRepo expected to be a singleton; actual: %+v", actual)
	}
}

func TestExampleServiceComponent(t *testing.T) {
	actual := r.ExampleService()
	if actual != r.exampleService {
		t.Errorf("resolver.exampleService expected to be a singleton; actual: %+v", actual)
	}
}

func TestHTTPServerComponent(t *testing.T) {
	actual := r.HTTPServer()
	if actual != r.httpServer {
		t.Errorf("resolver.httpServer expected to be a singleton; actual: %+v", actual)
	}
}

func TestLogComponent(t *testing.T) {
	actual := r.Log()
	if actual != r.log {
		t.Errorf("resolver.log expected to be a singleton; actual: %+v", actual)
	}
}

func TestMetadataComponent(t *testing.T) {
	actual := r.Metadata()
	if actual != r.metadata {
		t.Errorf("resolver.metadata expected to be a singleton; actual: %+v", actual)
	}
}

func TestPostgreSQLClientComponent(t *testing.T) {
	actual := r.PostgreSQLClient()
	if actual != r.postgreSQLClient {
		t.Errorf("resolver.postgreSQLClient expected to be a singleton; actual: %+v", actual)
	}
}
