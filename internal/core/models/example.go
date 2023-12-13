package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/jsonapi"
	"github.com/jasonsites/gosk/internal/core/pagination"
)

// ExampleRequestData defines the subset of Example domain model attributes that are accepted
// for input data request binding
type ExampleRequestData struct {
	Deleted     bool    `json:"deleted" validate:"omitempty,boolean"`
	Description *string `json:"description" validate:"omitempty,min=3,max=999"`
	Enabled     bool    `json:"enabled"  validate:"omitempty,boolean"`
	Status      *uint32 `json:"status" validate:"omitempty,numeric"`
	Title       string  `json:"title" validate:"required,omitempty,min=2,max=255"`
}

// ExampleDomainModel an Example domain model that contains one or more ExampleObject(s)
// and related metadata
type ExampleDomainModel struct {
	Data []ExampleObject
	Meta *ModelMetadata
	Solo bool
}

type ModelMetadata struct {
	Paging pagination.PageMetadata
}

// ExampleObject
type ExampleObject struct {
	Attributes ExampleObjectAttributes
	Meta       any
	Related    any
}

// Example defines an Example domain model for application logic
type ExampleObjectAttributes struct {
	ID          uuid.UUID  `json:"-"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      *uint32    `json:"status"`
	Enabled     bool       `json:"enabled"`
	Deleted     bool       `json:"-"`
	CreatedOn   time.Time  `json:"created_on"`
	CreatedBy   uint32     `json:"created_by"`
	ModifiedOn  *time.Time `json:"modified_on"`
	ModifiedBy  *uint32    `json:"modified_by"`
}

func (m *ExampleDomainModel) FormatDetailResponse() *jsonapi.Response[jsonapi.ResponseResource[ExampleObjectAttributes]] {
	resource := formatResource(&m.Data[0])
	response := &jsonapi.Response[jsonapi.ResponseResource[ExampleObjectAttributes]]{Data: resource}

	return response
}

func (m *ExampleDomainModel) FormatListResponse() *jsonapi.Response[[]jsonapi.ResponseResource[ExampleObjectAttributes]] {
	meta := &jsonapi.ResponseMetadata{
		Paging: pagination.PageMetadata{
			Limit:  m.Meta.Paging.Limit,
			Offset: m.Meta.Paging.Offset,
			Total:  m.Meta.Paging.Total,
		},
	}

	data := make([]jsonapi.ResponseResource[ExampleObjectAttributes], 0, len(m.Data))
	for _, domo := range m.Data {
		resource := formatResource(&domo)
		data = append(data, resource)
	}
	response := &jsonapi.Response[[]jsonapi.ResponseResource[ExampleObjectAttributes]]{
		Meta: meta,
		Data: data,
	}

	return response
}

// serializeResource
func formatResource(domo *ExampleObject) jsonapi.ResponseResource[ExampleObjectAttributes] {
	return jsonapi.ResponseResource[ExampleObjectAttributes]{
		Type: "example", // TODO
		ID:   domo.Attributes.ID,
		// Meta: domo.Meta,
		Attributes: ExampleObjectAttributes{
			Title:       domo.Attributes.Title,
			Description: domo.Attributes.Description,
			Status:      domo.Attributes.Status,
			Enabled:     domo.Attributes.Enabled,
			Deleted:     domo.Attributes.Deleted,
			CreatedOn:   domo.Attributes.CreatedOn,
			CreatedBy:   domo.Attributes.CreatedBy,
			ModifiedOn:  domo.Attributes.ModifiedOn,
			ModifiedBy:  domo.Attributes.ModifiedBy,
		},
	}
}
