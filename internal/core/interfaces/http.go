package interfaces

import (
	"net/http"

	"github.com/jasonsites/gosk/internal/core/jsonapi"
)

// ResourceController
type ResourceController interface {
	ResourceCreator
	ResourceDeleter
	ResourceDetailRetriever
	ResourceListRetriever
	ResourceUpdater
}

type ResourceCreator interface {
	Create(func() *jsonapi.RequestBody) http.HandlerFunc
}

type ResourceDeleter interface {
	Delete() http.HandlerFunc
}

type ResourceDetailRetriever interface {
	Detail() http.HandlerFunc
}

type ResourceListRetriever interface {
	List() http.HandlerFunc
}

type ResourceUpdater interface {
	Update(func() *jsonapi.RequestBody) http.HandlerFunc
}
