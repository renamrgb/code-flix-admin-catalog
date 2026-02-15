package create

import (
	"errors"
	"testing"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type CategoryGatewayMock struct {
	CreateFn func(*category.Category) (*category.Category, error)
}

func (m *CategoryGatewayMock) CreateCategory(cat *category.Category) (*category.Category, error) {
	return m.CreateFn(cat)
}

func (m *CategoryGatewayMock) GetCategoryByID(id category.CategoryID) (*category.Category, error) {
	return nil, nil
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

func TestCreateCategoryUseCaseExecute(t *testing.T) {
	gateway := &CategoryGatewayMock{
		CreateFn: func(cat *category.Category) (*category.Category, error) {
			return cat, nil
		},
	}

	useCase := NewCreateCategoryUseCase(gateway)

	input := CreateCategoryInput{
		Name:        "Movies",
		Description: "some description",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == nil {
		t.Fatal("expected output")
	}

	if output.ID == "" {
		t.Fatal("expected valid category ID")
	}
}

func TestCreateCategoryUseCase_ValidationError(t *testing.T) {
	gateway := &CategoryGatewayMock{
		CreateFn: func(cat *category.Category) (*category.Category, error) {
			return cat, nil
		},
	}

	useCase := NewCreateCategoryUseCase(gateway)

	input := CreateCategoryInput{
		Name:        "",
		Description: "some description",
		IsActive:    true,
	}

	output, err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected validation error")
	}

	if output != nil {
		t.Fatal("expected nil output on error")
	}
}

func TestCreateCategoryUseCase_GatewayError(t *testing.T) {
	expectedErr := errors.New("database error")

	gateway := &CategoryGatewayMock{
		CreateFn: func(cat *category.Category) (*category.Category, error) {
			return nil, expectedErr
		},
	}

	useCase := NewCreateCategoryUseCase(gateway)

	input := CreateCategoryInput{
		Name:        "Movies",
		Description: "some description",
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
		t.Fatal("expected nil output on error")
	}
}
