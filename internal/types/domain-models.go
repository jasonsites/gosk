package types

import (
	"encoding/json"
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
func (m *ExampleDomainModel) Serialize() ([]byte, error) {
	var result []byte

	// single resource response
	if m.Solo {
		resource := serializeResource(&m.Data[0])
		response := &JSONResponseSolo{Data: resource}
		result, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	// multiple resource response
	meta := &APIMetadata{
		Paging: PageMetadata{
			Limit:  m.Meta.Paging.Limit,
			Offset: m.Meta.Paging.Offset,
			Total:  m.Meta.Paging.Total,
		},
	}

	data := make([]ResponseResource, 0)
	for _, domo := range m.Data {
		resource := serializeResource(&domo)
		data = append(data, resource)
	}
	response := &JSONResponseMult{
		Meta: meta,
		Data: data,
	}
	result, err := json.Marshal(response)
	if err != nil {
		return nil, err
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
