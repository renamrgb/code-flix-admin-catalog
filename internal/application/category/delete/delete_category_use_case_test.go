package delete

import (
	"errors"
	"testing"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type CategoryGatewayMock struct {
	DeleteFn func(category.CategoryID) error
}

func (m *CategoryGatewayMock) CreateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryGatewayMock) GetCategoryByID(id category.CategoryID) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryGatewayMock) UpdateCategory(cat *category.Category) (*category.Category, error) {
	return nil, nil
}

func (m *CategoryGatewayMock) DeleteCategory(id category.CategoryID) error {
	return m.DeleteFn(id)
}

func (m *CategoryGatewayMock) FindAll(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
	return nil, nil
}

func TestDeleteCategoryUseCaseExecute(t *testing.T) {
	catID := category.NewCategoryID()

	var receivedID category.CategoryID

	gateway := &CategoryGatewayMock{
		DeleteFn: func(id category.CategoryID) error {
			receivedID = id
			return nil
		},
	}

	useCase := NewDeleteCategoryUseCase(gateway)

	input := DeleteCategoryInput{
		ID: catID.String(),
	}

	err := useCase.Execute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedID != catID {
		t.Fatal("expected correct category ID to be passed to gateway")
	}
}

func TestDeleteCategoryUseCase_InvalidID(t *testing.T) {
	gateway := &CategoryGatewayMock{
		DeleteFn: func(id category.CategoryID) error {
			return nil
		},
	}

	useCase := NewDeleteCategoryUseCase(gateway)

	input := DeleteCategoryInput{
		ID: "invalid-uuid", // 😈
	}

	err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected error for invalid UUID")
	}
}

func TestDeleteCategoryUseCase_GatewayError(t *testing.T) {
	expectedErr := errors.New("database error")

	catID := category.NewCategoryID()

	gateway := &CategoryGatewayMock{
		DeleteFn: func(id category.CategoryID) error {
			return expectedErr
		},
	}

	useCase := NewDeleteCategoryUseCase(gateway)

	input := DeleteCategoryInput{
		ID: catID.String(),
	}

	err := useCase.Execute(input)

	if err == nil {
		t.Fatal("expected gateway error")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: %v", err)
	}
}
