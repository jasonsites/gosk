package fixtures

import (
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/jasonsites/gosk/internal/core/models"
)

func ExampleModel(input *models.ExampleInputData) *models.ExampleInputData {

	description := fake.Sentence(4)
	status := uint32(fake.IntRange(0, 9))

	return &models.ExampleInputData{
		Deleted:     false,
		Description: &description,
		Enabled:     fake.Bool(),
		Status:      &status,
		Title:       fake.JobTitle(),
	}
}
