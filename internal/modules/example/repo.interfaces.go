package example

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	repo "github.com/jasonsites/gosk/internal/modules/common/repository"
)

// ExampleEntity defines an Example database entity
type ExampleEntity struct {
	ID              uuid.UUID
	Title           string
	Description     sql.NullString
	Status          repo.RecordStatus
	CreatedOn       time.Time
	CreatedContext  []byte // JSONB field
	ModifiedOn      time.Time
	ModifiedContext []byte // JSONB field
}

type ExampleEntityModel struct {
	Record ExampleEntity
}

// exampleEntityDefinition
type exampleEntityDefinition struct {
	Field exampleEntityFieldMap
	Name  string
}

// exampleEntityFields
type exampleEntityFieldMap struct {
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
	Field: exampleEntityFieldMap{
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
