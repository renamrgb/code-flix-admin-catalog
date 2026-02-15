// Package retrive provides use cases for retrieving category information by ID.
package retrive

import "github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"

type GetCategoryByIDUseCase struct {
	Gateway category.CategoryGateway
}

type GetCategoryByIDInput struct {
	ID string
}

func NewGetCategoryByIDUseCase(gateway category.CategoryGateway) *GetCategoryByIDUseCase {
	return &GetCategoryByIDUseCase{
		Gateway: gateway,
	}
}

func (uc *GetCategoryByIDUseCase) Execute(input GetCategoryByIDInput) (*category.Category, error) {
	id, err := category.ParseCategoryID(input.ID)
	if err != nil {
		return nil, err
	}

	return uc.Gateway.GetCategoryByID(id)
}
