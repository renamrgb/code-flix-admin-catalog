// Package update provides use cases for updating categories in the admin catalog.
package update

import "github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"

type UpdateCategoryUseCase struct {
	Gateway category.CategoryGateway
}

type UpdateCategoryInput struct {
	ID          string
	Name        string
	Description string
	IsActive    bool
}

type UpdateCategoryOutput struct {
	ID category.CategoryID
}

func NewUpdateCategoryUseCase(gateway category.CategoryGateway) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		Gateway: gateway,
	}
}

func (uc *UpdateCategoryUseCase) Execute(input UpdateCategoryInput) (*UpdateCategoryOutput, error) {
	id, err := category.ParseCategoryID(input.ID)
	if err != nil {
		return nil, err
	}
	cat, err := uc.Gateway.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	cat.Update(input.Name, input.Description, input.IsActive)

	if err := cat.Validate(); err != nil {
		return nil, err
	}

	cat, err = uc.Gateway.UpdateCategory(cat)
	if err != nil {
		return nil, err
	}

	return &UpdateCategoryOutput{
		ID: cat.ID,
	}, nil
}
