package types

import (
	"context"

	"github.com/google/uuid"
)

// RepoCreator
type RepoCreator interface {
	Create(context.Context, any) (DomainModel, error)
}

// RepoDeleter
type RepoDeleter interface {
	Delete(context.Context, uuid.UUID) error
}

// RepoDetailRetriever
type RepoDetailRetriever interface {
	Detail(context.Context, uuid.UUID) (DomainModel, error)
}

// RepoListRetriever
type RepoListRetriever interface {
	List(context.Context, QueryData) (DomainModel, error)
}

// RepoUpdater
type RepoUpdater interface {
	Update(context.Context, any, uuid.UUID) (DomainModel, error)
}

// ExampleRepository
type ExampleRepository interface {
	Create(context.Context, *ExampleRequestData) (DomainModel, error)
	RepoDeleter
	RepoDetailRetriever
	RepoListRetriever
	Update(context.Context, *ExampleRequestData, uuid.UUID) (DomainModel, error)
}

// temp documentation
// {
// 		meta: {
// 				paging: {
// 						limit,
// 						offset,
// 						total
// 				}
// 		},
// 		data: [{
// 				type: 'resource-type',
// 				meta: {
// 						...resource metadata
// 				},
// 				attributes: {
// 						...resource attributes
// 				},
// 				rel: [{
// 						type: 'rel-type',
// 						data: [{
// 								...rel-resource
// 						}],
// 				}],
// 		}]
// }
