package update

import (
	"errors"
	"testing"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type CategoryGatewayMock struct {
	GetByIDFn func(category.CategoryID) (*category.Category, error)
	UpdateFn  func(*category.Category) (*category.Category, error)
}

func (m *CategoryGatewayMock) CreateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryGatewayMock) GetCategoryByID(id category.CategoryID) (*category.Category, error) {
	return m.GetByIDFn(id)
}

func (m *CategoryGatewayMock) UpdateCategory(cat *category.Category) (*category.Category, error) {
	return m.UpdateFn(cat)
}

func (m *CategoryGatewayMock) DeleteCategory(id category.CategoryID) error {
	return nil
}

func (m *CategoryGatewayMock) FindAll(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
	return nil, nil
}

func TestUpdateCategoryUseCase_Execute(t *testing.T) {
	existingCategory, _ := category.NewCategory(
		"Movies",
		"old description",
		true,
	)

	gateway := &CategoryGatewayMock{
		GetByIDFn: func(id category.CategoryID) (*category.Category, error) {
			return existingCategory, nil
		},
		UpdateFn: func(cat *category.Category) (*category.Category, error) {
			return cat, nil
		},
	}

	useCase := NewUpdateCategoryUseCase(gateway)

	input := UpdateCategoryInput{
		ID:          existingCategory.ID.String(),
		Name:        "Updated Movies",
		Description: "new description",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == nil {
		t.Fatal("expected output")
	}

	if output.ID != existingCategory.ID {
		t.Fatal("expected same category ID")
	}

	if existingCategory.Name != input.Name {
		t.Errorf("expected name %s, got %s", input.Name, existingCategory.Name)
	}

	if existingCategory.Description != input.Description {
		t.Errorf("expected description %s, got %s", input.Description, existingCategory.Description)
	}
}

func TestUpdateCategoryUseCase_InvalidID(t *testing.T) {
	gateway := &CategoryGatewayMock{}

	useCase := NewUpdateCategoryUseCase(gateway)

	input := UpdateCategoryInput{
		ID:          "invalid-uuid", // 😈
		Name:        "Movies",
		Description: "desc",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected error for invalid ID")
	}

	if output != nil {
		t.Fatal("expected nil output")
	}
}

func TestUpdateCategoryUseCase_GetByIDError(t *testing.T) {
	expectedErr := errors.New("database error")

	existingID := category.NewCategoryID()

	gateway := &CategoryGatewayMock{
		GetByIDFn: func(id category.CategoryID) (*category.Category, error) {
			return nil, expectedErr
		},
	}

	useCase := NewUpdateCategoryUseCase(gateway)

	input := UpdateCategoryInput{
		ID:          existingID.String(),
		Name:        "Movies",
		Description: "desc",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected gateway error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != nil {
		t.Fatal("expected nil output")
	}
}

func TestUpdateCategoryUseCase_ValidationError(t *testing.T) {
	existingCategory, _ := category.NewCategory(
		"Movies",
		"description",
		true,
	)

	gateway := &CategoryGatewayMock{
		GetByIDFn: func(id category.CategoryID) (*category.Category, error) {
			return existingCategory, nil
		},
		UpdateFn: func(cat *category.Category) (*category.Category, error) {
			return cat, nil
		},
	}

	useCase := NewUpdateCategoryUseCase(gateway)

	input := UpdateCategoryInput{
		ID:          existingCategory.ID.String(),
		Name:        "", // inválido 😈
		Description: "desc",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected validation error")
	}

	if output != nil {
		t.Fatal("expected nil output")
	}
}

func TestUpdateCategoryUseCase_UpdateError(t *testing.T) {
	expectedErr := errors.New("update error")

	existingCategory, _ := category.NewCategory(
		"Movies",
		"description",
		true,
	)

	gateway := &CategoryGatewayMock{
		GetByIDFn: func(id category.CategoryID) (*category.Category, error) {
			return existingCategory, nil
		},
		UpdateFn: func(cat *category.Category) (*category.Category, error) {
			return nil, expectedErr
		},
	}

	useCase := NewUpdateCategoryUseCase(gateway)

	input := UpdateCategoryInput{
		ID:          existingCategory.ID.String(),
		Name:        "Updated Movies",
		Description: "desc",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected update error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: %v", err)
	}

	if output != nil {
		t.Fatal("expected nil output")
	}
}
