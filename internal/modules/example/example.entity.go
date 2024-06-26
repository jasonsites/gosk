package example

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/modules/common/pagination"
)

// ExampleEntity defines an Example database entity
type ExampleEntity struct {
	CreatedBy   uint32
	CreatedOn   time.Time
	Deleted     bool
	Description sql.NullString
	Enabled     bool
	ID          uuid.UUID
	ModifiedBy  sql.NullInt32
	ModifiedOn  sql.NullTime
	Status      sql.NullInt32
	Title       string
}

type ExampleEntityModel struct {
	Data []ExampleEntity
	Meta *ModelMetadata
	Solo bool
}

func (m *ExampleEntityModel) Unmarshal() *ExampleDomainModel {
	// single entity model
	if m.Solo {
		edo := unmarshalEntity(m.Data[0])
		model := &ExampleDomainModel{
			Data: []ExampleObject{*edo},
			Solo: m.Solo,
		}

		return model
	}

	// multiple entity model
	meta := &ModelMetadata{
		Paging: pagination.PageMetadata{
			Limit:  m.Meta.Paging.Limit,
			Offset: m.Meta.Paging.Offset,
			Total:  m.Meta.Paging.Total,
		},
	}

	data := make([]ExampleObject, 0, len(m.Data))
	// TODO: go routine?
	for _, record := range m.Data {
		edo := unmarshalEntity(record)
		data = append(data, *edo)
	}

	result := &ExampleDomainModel{
		Meta: meta,
		Data: data,
	}

	return result
}

func unmarshalEntity(e ExampleEntity) *ExampleObject {
	var (
		description *string
		modifiedBy  *uint32
		modifiedOn  *time.Time
		status      *uint32
	)

	if e.Description.Valid {
		description = &e.Description.String
	}
	if e.ModifiedBy.Valid {
		m := uint32(e.ModifiedBy.Int32)
		modifiedBy = &m
	}
	if e.ModifiedOn.Valid {
		modifiedOn = &e.ModifiedOn.Time
	}
	if e.Status.Valid {
		s := uint32(e.Status.Int32)
		status = &s
	}

	attributes := ExampleObjectAttributes{
		CreatedBy:   e.CreatedBy,
		CreatedOn:   e.CreatedOn,
		Deleted:     e.Deleted,
		Description: description,
		Enabled:     e.Enabled,
		ID:          e.ID,
		ModifiedBy:  modifiedBy,
		ModifiedOn:  modifiedOn,
		Status:      status,
		Title:       e.Title,
	}

	return &ExampleObject{
		Attributes: attributes,
	}
}
