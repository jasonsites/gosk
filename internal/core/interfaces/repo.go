package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk-api/internal/core/models"
	"github.com/jasonsites/gosk-api/internal/core/query"
)

// ModelCreator
type ModelCreator interface {
	Create(context.Context, any) (DomainModel, error)
}

// ModelDeleter
type ModelDeleter interface {
	Delete(context.Context, uuid.UUID) error
}

// ModelDetailRetriever
type ModelDetailRetriever interface {
	Detail(context.Context, uuid.UUID) (DomainModel, error)
}

// ModelListRetriever
type ModelListRetriever interface {
	List(context.Context, query.QueryData) (DomainModel, error)
}

// ModelUpdater
type ModelUpdater interface {
	Update(context.Context, any, uuid.UUID) (DomainModel, error)
}

// ExampleRepository defines the interface for repository managing the Example domain/entity model
type ExampleRepository interface {
	Create(context.Context, *models.ExampleRequestData) (DomainModel, error)
	ModelDeleter
	ModelDetailRetriever
	ModelListRetriever
	Update(context.Context, *models.ExampleRequestData, uuid.UUID) (DomainModel, error)
}
