package interfaces

import "github.com/jasonsites/gosk/internal/core/jsonapi"

// DomainModel defines the interface for all domain models
type DomainModel[T any] interface {
	ResponseFormatter[T]
}

type ResponseFormatter[T any] interface {
	FormatDetailResponse() (*jsonapi.Response[jsonapi.ResponseResource[T]], error)
	FormatListResponse() (*jsonapi.Response[[]jsonapi.ResponseResource[T]], error)
}
