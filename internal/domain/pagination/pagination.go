// Package pagination provides utilities for paginating collections of items.
package pagination

type Pagination[T any] struct {
	CurrentPage int
	PerPage     int
	Total       int
	Items       []T
}
