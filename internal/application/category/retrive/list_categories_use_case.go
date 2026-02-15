package retrive

import (
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type ListCategoriesUseCase struct {
	Gateway category.CategoryGateway
}

type ListCategoriesInput struct {
	Page      int
	PerPage   int
	Terms     string
	Sort      string
	Direction string
}

func NewListCategoriesUseCase(gateway category.CategoryGateway) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		Gateway: gateway,
	}
}

func (uc *ListCategoriesUseCase) Execute(input ListCategoriesInput) (*pagination.Pagination[category.Category], error) {
	query := category.SearchCategoryQuery{
		Page:      input.Page,
		PerPage:   input.PerPage,
		Terms:     input.Terms,
		Sort:      input.Sort,
		Direction: input.Direction,
	}

	return uc.Gateway.FindAll(query)
}
