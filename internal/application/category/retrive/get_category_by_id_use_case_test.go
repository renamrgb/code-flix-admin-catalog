// Package retrive provides use cases for retrieving category information by ID.
package retrive

import (
	"errors"
	"testing"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type CategoryGatewayMock struct {
	GetByIDFn func(category.CategoryID) (*category.Category, error)
}

func (m *CategoryGatewayMock) CreateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryGatewayMock) GetCategoryByID(id category.CategoryID) (*category.Category, error) {
	return m.GetByIDFn(id)
}

func (m *CategoryGatewayMock) UpdateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryGatewayMock) DeleteCategory(id category.CategoryID) error {
	return nil
}

func (m *CategoryGatewayMock) FindAll(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
	return nil, nil
}

func TestGetCategoryByIDUseCase_Execute(t *testing.T) {
	expectedCategory, _ := category.NewCategory(
		"Movies",
		"some description",
		true,
	)

	var receivedID category.CategoryID

	gateway := &CategoryGatewayMock{
		GetByIDFn: func(id category.CategoryID) (*category.Category, error) {
			receivedID = id
			return expectedCategory, nil
		},
	}

	useCase := NewGetCategoryByIDUseCase(gateway)

	input := GetCategoryByIDInput{
		ID: expectedCategory.ID.String(),
	}

	result, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected category")
	}

	if receivedID != expectedCategory.ID {
		t.Fatal("expected correct ID passed to gateway")
	}
}

func TestGetCategoryByIDUseCase_InvalidID(t *testing.T) {
	gateway := &CategoryGatewayMock{}

	useCase := NewGetCategoryByIDUseCase(gateway)

	input := GetCategoryByIDInput{
		ID: "invalid-uuid",
	}

	result, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected error")
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestGetCategoryByIDUseCase_GatewayError(t *testing.T) {
	expectedErr := errors.New("database error")

	catID := category.NewCategoryID()

	gateway := &CategoryGatewayMock{
		GetByIDFn: func(id category.CategoryID) (*category.Category, error) {
			return nil, expectedErr
		},
	}

	useCase := NewGetCategoryByIDUseCase(gateway)

	input := GetCategoryByIDInput{
		ID: catID.String(),
	}

	result, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected gateway error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Fatal("expected nil result")
	}
}
