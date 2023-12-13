package interfaces

import (
	"net/http"
)

// Controller
type Controller interface {
	Create() http.HandlerFunc
	Delete() http.HandlerFunc
	Detail() http.HandlerFunc
	List() http.HandlerFunc
	Update() http.HandlerFunc
}
