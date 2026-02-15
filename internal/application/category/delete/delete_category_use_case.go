// Package delete provides use cases for deleting categories in the application.
package delete

import "github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"

type DeleteCategoryUseCase struct {
	Gateway category.CategoryGateway
}

type DeleteCategoryInput struct {
	ID string
}

func NewDeleteCategoryUseCase(gateway category.CategoryGateway) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		Gateway: gateway,
	}
}

func (uc *DeleteCategoryUseCase) Execute(input DeleteCategoryInput) error {
	id, err := category.ParseCategoryID(input.ID)
	if err != nil {
		return err
	}

	return uc.Gateway.DeleteCategory(id)
}
