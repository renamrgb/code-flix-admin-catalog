package category

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/validation"
)

func TestNewCategory(t *testing.T) {
	tests := []struct {
		name          string
		isActive      bool
		expectDeleted bool
	}{
		{"active category", true, false},
		{"inactive category", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := "Movie"
			description := "some description"

			before := time.Now()

			cat, err := NewCategory(name, description, tt.isActive)

			after := time.Now()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cat == nil {
				t.Fatal("category should not be nil")
			}

			if cat.ID == (CategoryID{}) {
				t.Fatal("ID should be set")
			}

			if _, err := uuid.FromString(cat.ID.String()); err != nil {
				t.Fatalf("invalid UUID generated: %v", err)
			}

			if cat.IsActive != tt.isActive {
				t.Errorf("expected isActive %v, got %v", tt.isActive, cat.IsActive)
			}

			if cat.CreatedAt.IsZero() {
				t.Error("CreatedAt should be set")
			}

			if cat.UpdatedAt.IsZero() {
				t.Error("UpdatedAt should be set")
			}

			if tt.expectDeleted && cat.DeletedAt.IsZero() {
				t.Error("DeletedAt should be set for inactive category")
			}

			if !tt.expectDeleted && !cat.DeletedAt.IsZero() {
				t.Error("DeletedAt should be zero value for active category")
			}

			if cat.CreatedAt.Before(before) || cat.CreatedAt.After(after) {
				t.Error("CreatedAt timestamp is invalid")
			}

			if cat.UpdatedAt.Before(before) || cat.UpdatedAt.After(after) {
				t.Error("UpdatedAt timestamp is invalid")
			}
		})
	}
}

func TestCategoryUpdate(t *testing.T) {
	cat, _ := NewCategory("Old Name", "Old Desc", true)

	before := time.Now()

	cat.Update("New Name", "New Desc", false)

	after := time.Now()

	if cat.IsActive {
		t.Error("category should be inactive")
	}

	if cat.DeletedAt.IsZero() {
		t.Error("DeletedAt should be set when updating to inactive")
	}

	if cat.Name != "New Name" {
		t.Errorf("expected name %s, got %s", "New Name", cat.Name)
	}

	if cat.Description != "New Desc" {
		t.Errorf("expected description %s, got %s", "New Desc", cat.Description)
	}

	if cat.UpdatedAt.Before(before) || cat.UpdatedAt.After(after) {
		t.Error("UpdatedAt timestamp is invalid")
	}
}

func TestCategoryUpdate_ActivateFlow(t *testing.T) {
	cat, _ := NewCategory("Old Name", "Old Desc", false)

	if cat.DeletedAt.IsZero() {
		t.Fatal("category should start inactive/deleted")
	}

	before := time.Now()

	cat.Update("New Name", "New Desc", true)

	after := time.Now()

	if !cat.IsActive {
		t.Error("category should be active")
	}

	if !cat.DeletedAt.IsZero() {
		t.Error("DeletedAt should be cleared when activating")
	}

	if cat.UpdatedAt.Before(before) || cat.UpdatedAt.After(after) {
		t.Error("UpdatedAt timestamp is invalid")
	}
}

func TestCategoryValidate(t *testing.T) {
	tests := []struct {
		name           string
		inputName      string
		expectedErrors int
	}{
		{"empty name", "", 2},
		{"blank name", "   ", 2},
		{"name too short", "Go", 1},
		{"valid category", "Movies", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat, _ := NewCategory(tt.inputName, "desc", true)

			err := cat.Validate()

			if tt.expectedErrors == 0 && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedErrors > 0 && err == nil {
				t.Fatal("expected validation error")
			}

			if err != nil {
				var validationErr validation.ValidationErrors

				if !errors.As(err, &validationErr) {
					t.Fatalf("expected ValidationErrors, got %v", err)
				}

				if len(validationErr.Errs) != tt.expectedErrors {
					t.Fatalf("expected %d errors, got %d",
						tt.expectedErrors,
						len(validationErr.Errs),
					)
				}
			}
		})
	}
}

func TestCategoryActivate(t *testing.T) {
	cat, _ := NewCategory("Movies", "desc", false)

	if cat.DeletedAt.IsZero() {
		t.Fatal("category should start as deleted/inactive")
	}

	before := time.Now()

	cat.Activate()

	after := time.Now()

	if !cat.DeletedAt.IsZero() {
		t.Error("DeletedAt should be cleared on activate")
	}

	if !cat.IsActive {
		t.Error("category should be active")
	}

	if cat.UpdatedAt.Before(before) || cat.UpdatedAt.After(after) {
		t.Error("UpdatedAt timestamp is invalid")
	}
}

func TestCategoryDeactivate(t *testing.T) {
	cat, _ := NewCategory("Movies", "desc", true)

	if !cat.DeletedAt.IsZero() {
		t.Fatal("category should start as active")
	}

	before := time.Now()

	cat.Deactivate()

	after := time.Now()

	if cat.DeletedAt.IsZero() {
		t.Error("DeletedAt should be set on deactivate")
	}

	if cat.IsActive {
		t.Error("category should be inactive")
	}

	if cat.UpdatedAt.Before(before) || cat.UpdatedAt.After(after) {
		t.Error("UpdatedAt timestamp is invalid")
	}
}
