package example

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jasonsites/gosk/internal/app"
	cerror "github.com/jasonsites/gosk/internal/cerror"
	"github.com/jasonsites/gosk/internal/http/trace"
	"github.com/jasonsites/gosk/internal/logger"
	query "github.com/jasonsites/gosk/internal/modules/common/models/query"
	repo "github.com/jasonsites/gosk/internal/modules/common/repository"
)

// exampleEntityDefinition
type exampleEntityDefinition struct {
	Field exampleEntityFields
	Name  string
}

// exampleEntityFields
type exampleEntityFields struct {
	ID              string
	Title           string
	Description     string
	Status          string
	CreatedContext  string
	CreatedOn       string
	ModifiedContext string
	ModifiedOn      string
}

// exampleEntity
var exampleEntity = exampleEntityDefinition{
	Name: "example_entity",
	Field: exampleEntityFields{
		ID:              "id",
		Title:           "title",
		Description:     "description",
		Status:          "status",
		CreatedContext:  "created_context",
		CreatedOn:       "created_on",
		ModifiedContext: "modified_context",
		ModifiedOn:      "modified_on",
	},
}

// ExampleRepoConfig defines the input to NewExampleRepository
type ExampleRepoConfig struct {
	DBClient *pgxpool.Pool        `validate:"required"`
	Logger   *logger.CustomLogger `validate:"required"`
}

// exampleRepository
type exampleRepository struct {
	Entity exampleEntityDefinition
	db     *pgxpool.Pool
	logger *logger.CustomLogger
}

// NewExampleRepository returns a new exampleRepository instance
func NewExampleRepository(c *ExampleRepoConfig) (*exampleRepository, error) {
	if err := app.Validator.Validate.Struct(c); err != nil {
		return nil, err
	}

	repo := &exampleRepository{
		Entity: exampleEntity,
		db:     c.DBClient,
		logger: c.Logger,
	}

	return repo, nil
}

// Create
func (r *exampleRepository) Create(ctx context.Context, data *ExampleDTORequest) (*ModelContainer, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := r.logger.CreateContextLogger(traceID)

	// build sql query
	query := func() string {
		var (
			statement = "INSERT INTO %s %s VALUES %s RETURNING %s"
			field     = r.Entity.Field
			name      = r.Entity.Name
		)

		insertFields, values := repo.BuildInsertFieldsAndValues(
			field.Title,
			field.Description,
			field.CreatedContext,
		)

		returnFields := repo.BuildReturnFields(
			field.ID,
			field.Title,
			field.Description,
			field.Status,
			field.CreatedContext,
			field.CreatedOn,
			field.ModifiedContext,
			field.ModifiedOn,
		)

		return fmt.Sprintf(statement, name, insertFields, values, returnFields)
	}()

	// gather data from request, handling for nullable fields
	requestData := data

	var (
		description *string
		// Create a default created context - in a real app this would come from auth
		createdContext = map[string]any{
			"user_id": "system", // placeholder
		}
	)

	if requestData.Description != nil {
		description = requestData.Description
	}

	// Marshal context to JSON
	createdContextJSON, err := json.Marshal(createdContext)
	if err != nil {
		log.Error("Failed to marshal created context: " + err.Error())
		return nil, err
	}

	// create new entity for db row scan and execute query
	entity := ExampleEntity{}
	if err := r.db.QueryRow(
		ctx,
		query,
		requestData.Title,
		description,
		createdContextJSON,
	).Scan(
		&entity.ID,
		&entity.Title,
		&entity.Description,
		&entity.Status,
		&entity.CreatedContext,
		&entity.CreatedOn,
		&entity.ModifiedContext,
		&entity.ModifiedOn,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	em := &ExampleEntityModel{
		Record: entity,
	}

	result := MarshalEntityModel(em)

	return result, nil
}

// Delete
func (r *exampleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := r.logger.CreateContextLogger(traceID)

	// build sql query
	query := func() string {
		var (
			statement = "DELETE FROM %s WHERE id = ('%s'::uuid) RETURNING %s"
			field     = r.Entity.Field
			name      = r.Entity.Name
		)

		returnFields := repo.BuildReturnFields(field.ID)

		return fmt.Sprintf(statement, name, id, returnFields)
	}()

	// create new entity for db row scan and execute query
	entity := ExampleEntity{}
	if err := r.db.QueryRow(ctx, query).Scan(&entity.ID); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Detail
func (r *exampleRepository) Detail(ctx context.Context, id uuid.UUID) (*ModelContainer, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := r.logger.CreateContextLogger(traceID)

	// build sql query
	query := func() string {
		var (
			statement = "SELECT %s FROM %s WHERE id = ('%s'::uuid)"
			field     = r.Entity.Field
			name      = r.Entity.Name
		)

		returnFields := repo.BuildReturnFields(
			field.ID,
			field.Title,
			field.Description,
			field.Status,
			field.CreatedContext,
			field.CreatedOn,
			field.ModifiedContext,
			field.ModifiedOn,
		)

		return fmt.Sprintf(statement, returnFields, name, id)
	}()

	// create new entity for db row scan and execute query
	entity := ExampleEntity{}
	if err := r.db.QueryRow(ctx, query).Scan(
		&entity.ID,
		&entity.Title,
		&entity.Description,
		&entity.Status,
		&entity.CreatedContext,
		&entity.CreatedOn,
		&entity.ModifiedContext,
		&entity.ModifiedOn,
	); err != nil {
		log.Error(err.Error())
		err := cerror.NewNotFoundError(nil, fmt.Sprintf("unable to find %s with id '%s'", r.Entity.Name, id))
		return nil, err
	}

	em := &ExampleEntityModel{
		Record: entity,
	}

	result := MarshalEntityModel(em)

	return result, nil
}

// List
func (r *exampleRepository) List(ctx context.Context, q query.QueryData) (*ModelContainer, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := r.logger.CreateContextLogger(traceID)

	var (
		limit  = *q.Page.Limit
		offset = *q.Page.Offset
	)

	// build sql query
	query := func() string {
		var (
			statement = "SELECT %s FROM %s ORDER BY %s %s LIMIT %s OFFSET %s"
			field     = r.Entity.Field
			name      = r.Entity.Name
			limit     = fmt.Sprint(limit)
			offset    = fmt.Sprint(offset)
		)

		// Get the first sort pair from the query
		sortPairs := q.Sort.GetSortPairs()
		var orderField, orderDir string
		if len(sortPairs) > 0 {
			orderField = sortPairs[0].Field
			orderDir = string(sortPairs[0].Order)
		} else {
			// Fallback to default
			orderField = "modified_on"
			orderDir = "desc"
		}

		returnFields := repo.BuildReturnFields(
			field.ID,
			field.Title,
			field.Description,
			field.Status,
			field.CreatedContext,
			field.CreatedOn,
			field.ModifiedContext,
			field.ModifiedOn,
		)

		return fmt.Sprintf(statement, returnFields, name, orderField, orderDir, limit, offset)
	}()

	// execute query, returning rows
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	// create new entity model
	ems := make([]*ExampleEntityModel, 0, limit)

	// scan row data into new entities, appending to repo result
	for rows.Next() {
		entity := ExampleEntity{}

		if err := rows.Scan(
			&entity.ID,
			&entity.Title,
			&entity.Description,
			&entity.Status,
			&entity.CreatedContext,
			&entity.CreatedOn,
			&entity.ModifiedContext,
			&entity.ModifiedOn,
		); err != nil {
			log.Error(err.Error())
			return nil, err
		}

		em := &ExampleEntityModel{
			Record: entity,
		}

		ems = append(ems, em)
	}

	if err := rows.Err(); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// TODO: Investigate https://stackoverflow.com/questions/28888375/run-a-query-with-a-limit-offset-and-also-get-the-total-number-of-rows
	// query for total count
	var total int
	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.Entity.Name)
	if err := r.db.QueryRow(ctx, totalQuery).Scan(&total); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	gmd := ListQueryData{
		Page: repo.PageData{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
		Sort: q.Sort.ToSortMetadata(),
	}

	result := MarshalEntityModelList(ems, gmd)

	return result, nil
}

// Update
func (r *exampleRepository) Update(ctx context.Context, data *ExampleDTORequest, id uuid.UUID) (*ModelContainer, error) {
	traceID := trace.GetTraceIDFromContext(ctx)
	log := r.logger.CreateContextLogger(traceID)

	// build sql query
	query := func() string {
		var (
			statement = "UPDATE %s SET %s WHERE id = ('%s'::uuid) RETURNING %s"
			name      = r.Entity.Name
			field     = r.Entity.Field
		)

		values := repo.BuildUpdateValues(
			field.Title,
			field.Description,
			field.ModifiedContext,
			field.ModifiedOn,
		)

		returnFields := repo.BuildReturnFields(
			field.ID,
			field.Title,
			field.Description,
			field.Status,
			field.CreatedContext,
			field.CreatedOn,
			field.ModifiedContext,
			field.ModifiedOn,
		)

		return fmt.Sprintf(statement, name, values, id, returnFields)
	}()

	// gather data from request, handling for nullable fields
	requestData := data

	var (
		description *string
		modifiedOn  = time.Now()
		// Create a default modified context - in a real app this would come from auth
		modifiedContext = map[string]any{
			"user_id": "system", // placeholder
		}
	)

	if requestData.Description != nil {
		description = requestData.Description
	}

	// Marshal context to JSON
	modifiedContextJSON, err := json.Marshal(modifiedContext)
	if err != nil {
		log.Error("Failed to marshal modified context: " + err.Error())
		return nil, err
	}

	// create new entity for db row scan and execute query
	entity := ExampleEntity{}
	if err := r.db.QueryRow(
		ctx,
		query,
		requestData.Title,
		description,
		modifiedContextJSON,
		modifiedOn,
	).Scan(
		&entity.ID,
		&entity.Title,
		&entity.Description,
		&entity.Status,
		&entity.CreatedContext,
		&entity.CreatedOn,
		&entity.ModifiedContext,
		&entity.ModifiedOn,
	); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	em := &ExampleEntityModel{
		Record: entity,
	}

	result := MarshalEntityModel(em)

	return result, nil
}
