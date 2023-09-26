package fixtures

import (
	"bytes"
	"encoding/json"

	"github.com/jasonsites/gosk/internal/core/jsonapi"
	"github.com/jasonsites/gosk/internal/core/models"
)

func ComposeJSONBody(body jsonapi.RequestBody) *bytes.Buffer {
	b, _ := json.Marshal(body)
	return bytes.NewBuffer([]byte(b))
}

func ExampleRequest(model *models.ExampleInputData) jsonapi.RequestBody {
	return jsonapi.RequestBody{
		Data: &jsonapi.RequestResource{
			Type:       "example",
			Attributes: model,
		}}
}
