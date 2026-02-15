package retrive

import (
	"errors"
	"testing"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type CategoryListGatewayMock struct {
	FindAllFn func(category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error)
}

func (m *CategoryListGatewayMock) CreateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryListGatewayMock) GetCategoryByID(id category.CategoryID) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryListGatewayMock) UpdateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryListGatewayMock) DeleteCategory(id category.CategoryID) error {
	return nil
}

func (m *CategoryListGatewayMock) FindAll(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
	return m.FindAllFn(query)
}

func TestListCategoriesUseCase_Execute(t *testing.T) {
	expectedCategories := []category.Category{
		{Name: "Movies"},
		{Name: "Series"},
	}

	expectedPagination := &pagination.Pagination[category.Category]{
		CurrentPage: 1,
		PerPage:     10,
		Total:       2,
		Items:       expectedCategories,
	}

	var receivedQuery category.SearchCategoryQuery

	gateway := &CategoryListGatewayMock{
		FindAllFn: func(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
			receivedQuery = query
			return expectedPagination, nil
		},
	}

	useCase := NewListCategoriesUseCase(gateway)

	input := ListCategoriesInput{
		Page:      1,
		PerPage:   10,
		Terms:     "mov",
		Sort:      "name",
		Direction: "asc",
	}

	result, err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected pagination result")
	}

	if result.Total != expectedPagination.Total {
		t.Fatalf("expected total %d, got %d", expectedPagination.Total, result.Total)
	}

	if receivedQuery.Page != input.Page {
		t.Error("invalid page mapping")
	}

	if receivedQuery.PerPage != input.PerPage {
		t.Error("invalid perPage mapping")
	}

	if receivedQuery.Terms != input.Terms {
		t.Error("invalid terms mapping")
	}

	if receivedQuery.Sort != input.Sort {
		t.Error("invalid sort mapping")
	}

	if receivedQuery.Direction != input.Direction {
		t.Error("invalid direction mapping")
	}
}

func TestListCategoriesUseCase_GatewayError(t *testing.T) {
	expectedErr := errors.New("database error")

	gateway := &CategoryListGatewayMock{
		FindAllFn: func(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
			return nil, expectedErr
		},
	}

	useCase := NewListCategoriesUseCase(gateway)

	input := ListCategoriesInput{
		Page:    1,
		PerPage: 10,
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
