package domain

import (
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/jasonsites/gosk-api/internal/validation"
)

// Domain is the top-level container for the application domain layer
type Domain struct {
	Services *Services
}

// Services contains all individual resource services
type Services struct {
	Example types.Service
}

// NewDomain creates a new Domain instance
func NewDomain(s *Services) (*Domain, error) {
	if err := validation.Validate.Struct(s); err != nil {
		return nil, err
	}

	return &Domain{Services: s}, nil
}
