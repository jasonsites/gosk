package fixtures

import (
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/jasonsites/gosk/internal/modules/example"
)

func ExampleModel(input *example.ExampleDTO) *example.ExampleDTO {

	description := fake.Sentence(4)
	status := uint32(fake.IntRange(0, 9))

	return &example.ExampleDTO{
		Deleted:     false,
		Description: &description,
		Enabled:     fake.Bool(),
		Status:      &status,
		Title:       fake.JobTitle(),
	}
}
