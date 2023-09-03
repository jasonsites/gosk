package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/jasonsites/gosk-api/internal/validation"
)

// ExampleServiceConfig
type ExampleServiceConfig struct {
	Logger *types.Logger           `validate:"required"`
	Repo   types.ExampleRepository `validate:"required"`
}

// exampleService
type exampleService struct {
	logger *types.Logger
	repo   types.ExampleRepository
}

// NewExampleService
func NewExampleService(c *ExampleServiceConfig) (*exampleService, error) {
	if err := validation.Validate.Struct(c); err != nil {
		return nil, err
	}

	log := c.Logger.Log.With().Str("tags", "service,resource").Logger()
	logger := &types.Logger{
		Enabled: c.Logger.Enabled,
		Level:   c.Logger.Level,
		Log:     &log,
	}

	service := &exampleService{
		logger: logger,
		repo:   c.Repo,
	}

	return service, nil
}

// Create
func (s *exampleService) Create(ctx context.Context, data any) (types.DomainModel, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	fmt.Printf("data in CREATE: %+v\n", data)

	d, ok := data.(*types.ExampleRequestData)
	if !ok {
		err := fmt.Errorf("error asserting data as types.ExampleRequestData")
		log.Error().Err(err).Send()
		return nil, err
	}

	model, err := s.repo.Create(ctx, d)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return model, nil
}

// Delete
func (s *exampleService) Delete(ctx context.Context, id uuid.UUID) error {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Error().Err(err).Send()
		return err
	}

	return nil
}

// Detail
func (s *exampleService) Detail(ctx context.Context, id uuid.UUID) (types.DomainModel, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	model, err := s.repo.Detail(ctx, id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return model, nil
}

// List
func (s *exampleService) List(ctx context.Context, q types.QueryData) (types.DomainModel, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	model, err := s.repo.List(ctx, q)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return model, nil
}

// Update
func (s *exampleService) Update(ctx context.Context, data any, id uuid.UUID) (types.DomainModel, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	model, err := s.repo.Update(ctx, data.(*types.ExampleRequestData), id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return model, nil
}
