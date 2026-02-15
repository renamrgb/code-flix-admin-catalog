// Package create provides use cases for creating categories in the admin catalog.
package create

import (
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
)

type CreateCategoryUseCase struct {
	Gateway category.CategoryGateway
}

type CreateCategoryInput struct {
	Name        string
	Description string
	IsActive    bool
}

type CreateCategoryOutput struct {
	ID string
}

func NewCreateCategoryUseCase(gateway category.CategoryGateway) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		Gateway: gateway,
	}
}

func (uc *CreateCategoryUseCase) Execute(input CreateCategoryInput) (*CreateCategoryOutput, error) {
	cat, err := category.NewCategory(
		input.Name,
		input.Description,
		input.IsActive,
	)
	if err != nil {
		return nil, err
	}

	if err := cat.Validate(); err != nil {
		return nil, err
	}

	cat, err = uc.Gateway.CreateCategory(cat)
	if err != nil {
		return nil, err
	}

	return &CreateCategoryOutput{
		ID: cat.ID.String(),
	}, nil
}
