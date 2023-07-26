package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/jasonsites/gosk-api/internal/validation"
)

// ExampleServiceConfig
type ExampleServiceConfig struct {
	Logger *types.Logger    `validate:"required"`
	Repo   types.Repository `validate:"required"`
}

// exampleService
type exampleService struct {
	logger *types.Logger
	repo   types.Repository
}

// NewExampleService
func NewExampleService(c *ExampleServiceConfig) (*exampleService, error) {
	if err := validation.Validate.Struct(c); err != nil {
		return nil, err
	}

	log := c.Logger.Log.With().Str("tags", "service,example").Logger()
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
func (s *exampleService) Create(ctx context.Context, data any) (*types.JSONResponseSolo, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	result, err := s.repo.Create(ctx, data.(*types.ExampleRequestData))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	model := &types.Example{}
	sr, err := model.SerializeResponse(result, true)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	res := sr.(*types.JSONResponseSolo)

	return res, nil
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
func (s *exampleService) Detail(ctx context.Context, id uuid.UUID) (*types.JSONResponseSolo, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	result, err := s.repo.Detail(ctx, id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	model := &types.Example{}
	sr, err := model.SerializeResponse(result, true)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	res := sr.(*types.JSONResponseSolo)

	return res, nil
}

// List
func (s *exampleService) List(ctx context.Context, q types.QueryData) (*types.JSONResponseMult, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	result, err := s.repo.List(ctx, q)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	model := &types.Example{}
	sr, err := model.SerializeResponse(result, false)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	res := sr.(*types.JSONResponseMult)

	return res, nil
}

// Update
func (s *exampleService) Update(ctx context.Context, data any, id uuid.UUID) (*types.JSONResponseSolo, error) {
	requestId := ctx.Value(types.CorrelationContextKey).(*types.Trace).RequestID
	log := s.logger.Log.With().Str("req_id", requestId).Logger()

	result, err := s.repo.Update(ctx, data.(*types.ExampleRequestData), id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	model := &types.Example{}
	sr, err := model.SerializeResponse(result, true)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}
	res := sr.(*types.JSONResponseSolo)

	return res, nil
}
