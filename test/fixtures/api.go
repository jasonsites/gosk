package fixtures

import (
	"bytes"
	"encoding/json"

	"github.com/jasonsites/gosk/internal/http/jsonapi"
	"github.com/jasonsites/gosk/internal/modules/example"
)

func ComposeJSONBody(body jsonapi.RequestBody) *bytes.Buffer {
	b, _ := json.Marshal(body)
	return bytes.NewBuffer([]byte(b))
}

func ExampleRequest(model *example.ExampleDTO) jsonapi.RequestBody {
	return jsonapi.RequestBody{
		Data: &jsonapi.RequestResource{
			Type:       "example",
			Attributes: model,
		}}
}
