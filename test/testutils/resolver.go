package testutils

import (
	"context"

	"github.com/jasonsites/gosk/internal/resolver"
)

// InitializeResolver creates a new Resolver from the given config and returns a reference to the configured Resolver
func InitializeResolver(conf *resolver.Config, entry string) (*resolver.Resolver, error) {
	r := resolver.NewResolver(context.Background(), conf)
	if entry == "" {
		entry = resolver.LoadEntries.HTTPServer
	}
	r.Load(entry)
	return r, nil
}
