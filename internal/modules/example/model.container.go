package example

import (
	"time"

	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/http/jsonapi"
	query "github.com/jasonsites/gosk/internal/modules/common/models/query"
)

// ModelContainer contains one or more ExampleModel(s) and related metadata
type ModelContainer struct {
	Data []ExampleModel
	Meta *ModelContainerMeta
	Solo bool
}

type ModelContainerMeta struct {
	Filter *query.FilterMetadata `json:"filter,omitempty"`
	Page   query.PageMetadata    `json:"page,omitempty"`
	Sort   query.SortMetadata    `json:"sort,omitempty"`
}

// ExampleModel
type ExampleModel struct {
	Meta       any
	Attributes ModelAttributes
}

// Example defines an Example domain model for application logic
type ModelAttributes struct {
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

func (m *ModelContainer) FormatResponse() (*jsonapi.Response, error) {
	if m.Solo {
		resource := formatResource(&m.Data[0])
		response := &jsonapi.Response{Data: resource}
		return response, nil
	}

	meta := &jsonapi.ResponseMetadata{
		Filter: m.Meta.Filter,
		Page: query.PageMetadata{
			Limit:  m.Meta.Page.Limit,
			Offset: m.Meta.Page.Offset,
			Total:  m.Meta.Page.Total,
		},
		Sort: &m.Meta.Sort,
	}

	data := make([]jsonapi.ResponseResource, 0, len(m.Data))
	for _, domo := range m.Data {
		resource := formatResource(&domo)
		data = append(data, resource)
	}
	response := &jsonapi.Response{
		Meta: meta,
		Data: data,
	}

	return response, nil
}

// serializeResource
func formatResource(domo *ExampleModel) jsonapi.ResponseResource {
	return jsonapi.ResponseResource{
		Type: "example", // TODO
		ID:   domo.Attributes.ID,
		// Meta: domo.Meta,
		Attributes: ModelAttributes{
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
