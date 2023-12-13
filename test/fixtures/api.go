package fixtures

import (
	"bytes"
	"encoding/json"

	"github.com/jasonsites/gosk/internal/core/jsonapi"
	"github.com/jasonsites/gosk/internal/core/models"
)

func ComposeJSONBody[T any](body jsonapi.RequestBody[T]) *bytes.Buffer {
	b, _ := json.Marshal(body)
	return bytes.NewBuffer([]byte(b))
}

func ExampleRequest(model *models.ExampleRequestData) jsonapi.RequestBody[models.ExampleRequestData] {
	return jsonapi.RequestBody[models.ExampleRequestData]{
		Data: &jsonapi.RequestResource[models.ExampleRequestData]{
			Type:       "example",
			Attributes: model,
		}}
}
