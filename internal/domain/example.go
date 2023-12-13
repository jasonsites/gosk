package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/interfaces"
	"github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/models"
	"github.com/jasonsites/gosk/internal/core/query"
	"github.com/jasonsites/gosk/internal/core/trace"
	"github.com/jasonsites/gosk/internal/core/validation"
)

// ExampleServiceConfig defines the input to NewExampleService
type ExampleServiceConfig struct {
	Logger *logger.CustomLogger         `validate:"required"`
	Repo   interfaces.ExampleRepository `validate:"required"`
}

// exampleService
type exampleService struct {
	logger *logger.CustomLogger
	repo   interfaces.ExampleRepository
}

// NewExampleService returns a new exampleService instance
func NewExampleService(c *ExampleServiceConfig) (*exampleService, error) {
	if err := validation.Validate.Struct(c); err != nil {
		return nil, err
	}

	service := &exampleService{
		logger: c.Logger,
		repo:   c.Repo,
	}

	return service, nil
}

// Create
func (s *exampleService) Create(ctx context.Context, data any) (*models.ExampleDomainModel, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := s.logger.CreateContextLogger(traceID)

	d, ok := data.(*models.ExampleRequestData)
	if !ok {
		err := fmt.Errorf("example input data assertion error")
		log.Error(err.Error())
		return nil, err
	}

	model, err := s.repo.Create(ctx, d)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return model, nil
}

// Delete
func (s *exampleService) Delete(ctx context.Context, id uuid.UUID) error {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := s.logger.CreateContextLogger(traceID)

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Detail
func (s *exampleService) Detail(ctx context.Context, id uuid.UUID) (*models.ExampleDomainModel, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := s.logger.CreateContextLogger(traceID)

	model, err := s.repo.Detail(ctx, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return model, nil
}

// List
func (s *exampleService) List(ctx context.Context, q query.QueryData) (*models.ExampleDomainModel, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := s.logger.CreateContextLogger(traceID)

	model, err := s.repo.List(ctx, q)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return model, nil
}

// Update
func (s *exampleService) Update(ctx context.Context, data any, id uuid.UUID) (*models.ExampleDomainModel, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := s.logger.CreateContextLogger(traceID)

	d, ok := data.(*models.ExampleRequestData)
	if !ok {
		err := fmt.Errorf("example input data assertion error")
		log.Error(err.Error())
		return nil, err
	}

	model, err := s.repo.Update(ctx, d, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return model, nil
}
