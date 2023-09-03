package types

import (
	"time"

	"github.com/google/uuid"
)

type ExampleDomainModel struct {
	Data []ExampleDomainObject
	Meta *ModelMetadata
	Solo bool
}

type ExampleDomainObject struct {
	Attributes ExampleDomainObjectAttributes
	Meta       any
	Related    any
}

// Example defines an Example domain model for application logic
type ExampleDomainObjectAttributes struct {
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

// SerializeResponse
func (m *ExampleDomainModel) Serialize() (JSONResponse, error) {
	// single resource response
	if m.Solo {
		resource := serializeResource(&m.Data[0])
		result := &JSONResponseSolo{Data: resource}

		return result, nil
	}

	// multiple resource response
	meta := &APIMetadata{
		Paging: ListPaging{
			Limit:  m.Meta.Paging.Limit,
			Offset: m.Meta.Paging.Offset,
			Total:  m.Meta.Paging.Total,
		},
	}

	data := make([]ResponseResource, 0)
	// TODO: go routine
	for _, domo := range m.Data {
		resource := serializeResource(&domo)
		data = append(data, resource)
	}

	result := &JSONResponseMult{
		Meta: meta,
		Data: data,
	}

	return result, nil
}

func serializeResource(domo *ExampleDomainObject) ResponseResource {
	return ResponseResource{
		Type: "example",
		ID:   domo.Attributes.ID,
		// Meta: domo.Meta,
		Attributes: ExampleDomainObjectAttributes{
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
