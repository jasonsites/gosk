package fixtures

import (
	"database/sql"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/entities"
)

func ExampleEntityRecord(input *entities.ExampleEntity, id *uuid.UUID) *entities.ExampleEntity {
	createdBy := uint32(fake.IntRange(1, 9999))
	createdOn := fake.Date()
	deleted := fake.Bool()
	description := sql.NullString{
		String: fake.Sentence(4),
		Valid:  true,
	}
	enabled := fake.Bool()
	modifiedBy := sql.NullInt32{
		Int32: int32(fake.IntRange(1, 9999)),
		Valid: true,
	}
	modifiedOn := sql.NullTime{
		Time:  fake.Date(),
		Valid: true,
	}
	status := sql.NullInt32{
		Int32: int32(fake.IntRange(1, 9)),
		Valid: true,
	}
	title := fake.JobTitle()

	if input != nil {
		deleted = input.Deleted
		enabled = input.Enabled

		if input.CreatedBy != 0 {
			createdBy = input.CreatedBy
		}
		if !input.CreatedOn.IsZero() {
			createdOn = input.CreatedOn
		}
		if input.Description.String != "" {
			description = input.Description
		}
		if input.ModifiedBy.Int32 != 0 {
			modifiedBy = input.ModifiedBy
		}
		if !input.ModifiedOn.Time.IsZero() {
			modifiedOn = input.ModifiedOn
		}
		if input.Status.Int32 != 0 {
			status = input.Status
		}
		if input.Title != "" {
			title = input.Title
		}
	}

	record := &entities.ExampleEntity{
		CreatedBy:   createdBy,
		CreatedOn:   createdOn,
		Deleted:     deleted,
		Description: description,
		Enabled:     enabled,
		ModifiedBy:  modifiedBy,
		ModifiedOn:  modifiedOn,
		Status:      status,
		Title:       title,
	}

	if id != nil {
		record.ID = *id
	}

	return record
}
