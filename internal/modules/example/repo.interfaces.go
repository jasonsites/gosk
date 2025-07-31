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
