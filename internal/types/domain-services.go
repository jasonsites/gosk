package types

import (
	"context"

	"github.com/google/uuid"
)

// DomainModel defines the interface for all domain models
type DomainModel interface {
	ResponseSerializer
}

// DomainRegistry defines a registry for all domain types to be used across the application
type DomainRegistry struct {
	Example string
}

// DomainType exposes constants for all domain types
var DomainType = DomainRegistry{
	Example: "example",
}

// TODO
// Discoverable defines the interface for all types with self discovery
type Discoverable interface {
	Discover() Discoverable
}

// TODO
// ResponseSerializer defines the interface for all types that serialize to JSON response
type ResponseSerializer interface {
	Serialize() (JSONResponse, error)
}

// Service
type Service interface {
	ServiceCreator
	ServiceDeleter
	ServiceDetailRetriever
	ServiceListRetriever
	ServiceUpdater
}

type ServiceCreator interface {
	Create(context.Context, any) (DomainModel, error)
}

type ServiceDeleter interface {
	Delete(context.Context, uuid.UUID) error
}

type ServiceDetailRetriever interface {
	Detail(context.Context, uuid.UUID) (DomainModel, error)
}

type ServiceListRetriever interface {
	List(context.Context, QueryData) (DomainModel, error)
}

type ServiceUpdater interface {
	Update(context.Context, any, uuid.UUID) (DomainModel, error)
}
