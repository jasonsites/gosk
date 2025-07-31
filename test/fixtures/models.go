package fixtures

import (
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/jasonsites/gosk/internal/modules/example"
)

func ExampleModel(input *example.ExampleDTORequest) *example.ExampleDTORequest {

	description := fake.Sentence(4)

	return &example.ExampleDTORequest{
		Description: &description,
		Title:       fake.JobTitle(),
	}
}
