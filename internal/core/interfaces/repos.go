package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/models"
	"github.com/jasonsites/gosk/internal/core/query"
)

// ExampleRepository defines the interface for repository managing the Example domain/entity model
type ExampleRepository interface {
	Create(context.Context, *models.ExampleDTO) (*models.ExampleDomainModel, error)
	Delete(context.Context, uuid.UUID) error
	Detail(context.Context, uuid.UUID) (*models.ExampleDomainModel, error)
	List(context.Context, query.QueryData) (*models.ExampleDomainModel, error)
	Update(context.Context, *models.ExampleDTO, uuid.UUID) (*models.ExampleDomainModel, error)
}
