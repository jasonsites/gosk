package example

import (
	"encoding/json"
	"time"

	query "github.com/jasonsites/gosk/internal/modules/common/models/query"
	repo "github.com/jasonsites/gosk/internal/modules/common/repository"
)

func MarshalEntityModel(em *ExampleEntityModel) *ModelContainer {
	model := marshalEntity(em.Record)

	container := &ModelContainer{
		Data: []ExampleModel{*model},
		Solo: true,
	}

	return container
}

func MarshalEntityModelList(ems []*ExampleEntityModel, lqd ListQueryData) *ModelContainer {
	meta := MarshalListMetadata(lqd)
	data := make([]ExampleModel, 0, len(ems))

	for _, em := range ems {
		edo := marshalEntity(em.Record)
		data = append(data, *edo)
	}

	result := &ModelContainer{
		Meta: meta,
		Data: data,
	}

	return result
}

type ListQueryData struct {
	Page repo.PageData
	Sort ExampleSortMetadata
}

func MarshalListMetadata(lqd ListQueryData) *ModelContainerMeta {
	meta := &ModelContainerMeta{
		Page: query.PageMetadata{
			Limit:  uint32(lqd.Page.Limit),
			Offset: uint32(lqd.Page.Offset),
			Total:  uint32(lqd.Page.Total),
		},
		Sort: lqd.Sort,
	}

	return meta
}

func marshalEntity(e ExampleEntity) *ExampleModel {
	var (
		description *string
		modifiedOn  *time.Time
	)

	if e.Description.Valid {
		description = &e.Description.String
	}

	// Since modified_on is NOT NULL in the schema, we don't need to check for validity
	modifiedOn = &e.ModifiedOn

	// Parse context data to extract user information for backward compatibility
	var createdBy uint32 = 0
	var modifiedBy *uint32

	// Try to parse created context
	if len(e.CreatedContext) > 0 {
		var createdCtx map[string]any
		if err := json.Unmarshal(e.CreatedContext, &createdCtx); err == nil {
			if userID, ok := createdCtx["user_id"]; ok {
				// This is a simplified conversion - in reality you'd handle this more robustly
				if userIDStr, ok := userID.(string); ok && userIDStr != "system" {
					// If it's a real user ID, you might convert it or look it up
					createdBy = 1 // placeholder conversion
				}
			}
		}
	}

	// Try to parse modified context
	if len(e.ModifiedContext) > 0 {
		var modifiedCtx map[string]any
		if err := json.Unmarshal(e.ModifiedContext, &modifiedCtx); err == nil {
			if userID, ok := modifiedCtx["user_id"]; ok {
				if userIDStr, ok := userID.(string); ok && userIDStr != "system" {
					modifiedByVal := uint32(1) // placeholder conversion
					modifiedBy = &modifiedByVal
				}
			}
		}
	}

	// Convert status enum to legacy format for backward compatibility
	var status *uint32
	switch e.Status {
	case repo.RecordStatusActive:
		statusVal := uint32(1)
		status = &statusVal
	case repo.RecordStatusArchived:
		statusVal := uint32(2)
		status = &statusVal
	case repo.RecordStatusDeleted:
		statusVal := uint32(3)
		status = &statusVal
	}

	// For backward compatibility, map status to enabled/deleted flags
	enabled := e.Status == repo.RecordStatusActive
	deleted := e.Status == repo.RecordStatusDeleted

	attributes := ModelAttributes{
		CreatedBy:   createdBy,
		CreatedOn:   e.CreatedOn,
		Deleted:     deleted,
		Description: description,
		Enabled:     enabled,
		ID:          e.ID,
		ModifiedBy:  modifiedBy,
		ModifiedOn:  modifiedOn,
		Status:      status,
		Title:       e.Title,
	}

	return &ExampleModel{
		Attributes: attributes,
	}
}
