package fixtures

import (
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/jasonsites/gosk/internal/core/models"
)

type Recipients struct {
	Recipients []string
	Count      int
}

func ExampleModel() *models.ExampleObjectAttributes {

	description := fake.Sentence(4)
	status := fake.Uint32()

	return &models.ExampleObjectAttributes{
		Description: &description,
		Enabled:     fake.Bool(),
		Status:      &status,
		Title:       fake.JobTitle(),
	}
}
