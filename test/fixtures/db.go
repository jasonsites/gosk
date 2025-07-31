package fixtures

import (
	"database/sql"
	"encoding/json"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	repo "github.com/jasonsites/gosk/internal/modules/common/repository"
	"github.com/jasonsites/gosk/internal/modules/example"
)

func ExampleEntityRecord(input *example.ExampleEntity, id *uuid.UUID) *example.ExampleEntity {
	createdOn := fake.Date()
	description := sql.NullString{
		String: fake.Sentence(4),
		Valid:  true,
	}
	modifiedOn := fake.Date()
	status := repo.RecordStatusActive
	title := fake.JobTitle()

	// Create default contexts
	createdContext := map[string]any{
		"user_id": "test_user",
	}
	modifiedContext := map[string]any{
		"user_id": "test_user",
	}

	createdContextJSON, _ := json.Marshal(createdContext)
	modifiedContextJSON, _ := json.Marshal(modifiedContext)

	if input != nil {
		if !input.CreatedOn.IsZero() {
			createdOn = input.CreatedOn
		}
		if input.Description.String != "" {
			description = input.Description
		}
		if !input.ModifiedOn.IsZero() {
			modifiedOn = input.ModifiedOn
		}
		if input.Status != "" {
			status = input.Status
		}
		if input.Title != "" {
			title = input.Title
		}
		if len(input.CreatedContext) > 0 {
			createdContextJSON = input.CreatedContext
		}
		if len(input.ModifiedContext) > 0 {
			modifiedContextJSON = input.ModifiedContext
		}
	}

	record := &example.ExampleEntity{
		CreatedOn:       createdOn,
		CreatedContext:  createdContextJSON,
		Description:     description,
		ModifiedOn:      modifiedOn,
		ModifiedContext: modifiedContextJSON,
		Status:          status,
		Title:           title,
	}

	if id != nil {
		record.ID = *id
	}

	return record
}
