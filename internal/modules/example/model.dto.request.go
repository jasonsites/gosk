package example

import (
	v "github.com/invopop/validation"
)

// ExampleDTORequest defines the subset of Example domain model attributes that are accepted
// for input data request binding
type ExampleDTORequest struct {
	Description *string `json:"description" validate:"omitempty,min=3,max=999"`
	Title       string  `json:"title" validate:"required,omitempty,min=2,max=255"`
}

// Validate validates an Example request DTO
func (e ExampleDTORequest) Validate() error {
	if err := v.ValidateStruct(&e,
		v.Field(&e.Title, v.Required),
	); err != nil {
		return err
	}

	return nil
}
