package category

import "github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"

type CategoryGateway interface {
	CreateCategory(category *Category) (*Category, error)
	GetCategoryByID(id CategoryID) (*Category, error)
	UpdateCategory(category *Category) (*Category, error)
	DeleteCategory(id CategoryID) error

	FindAll(query SearchCategoryQuery) (*pagination.Pagination[Category], error)
}

type SearchCategoryQuery struct {
	Page      int
	PerPage   int
	Terms     string
	Sort      string
	Direction string
}
