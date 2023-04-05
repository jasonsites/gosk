package domain

import (
	"github.com/jasonsites/gosk-api/internal/types"
)

// Domain is the top-level container for the application domain layer
type Domain struct {
	Services *Services
}

// Services contains all individual resource services
type Services struct {
	ResourceService types.Service
}

// NewDomain creates a new Domain instance
func NewDomain(s *Services) *Domain {
	return &Domain{Services: s}
}
