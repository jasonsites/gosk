package interfaces

import (
	"net/http"

	"github.com/jasonsites/gosk/internal/core/jsonapi"
)

// ExampleController
type ExampleController interface {
	Create(func() *jsonapi.RequestBody) http.HandlerFunc
	Delete() http.HandlerFunc
	Detail() http.HandlerFunc
	List() http.HandlerFunc
	Update(func() *jsonapi.RequestBody) http.HandlerFunc
}
