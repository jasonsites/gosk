package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/models"
	"github.com/jasonsites/gosk/internal/core/query"
)

// ExampleService
type ExampleService interface {
	Create(context.Context, any) (*models.ExampleDomainModel, error)
	Delete(context.Context, uuid.UUID) error
	Detail(context.Context, uuid.UUID) (*models.ExampleDomainModel, error)
	List(context.Context, query.QueryData) (*models.ExampleDomainModel, error)
	Update(context.Context, any, uuid.UUID) (*models.ExampleDomainModel, error)
}
