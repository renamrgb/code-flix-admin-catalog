// Package category provides domain logic for managing categories in the admin catalog.
package category

import (
	"errors"
	"strings"
	"time"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/validation"
)

type Category struct {
	ID          CategoryID
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func NewCategory(name, description string, isActive bool) (*Category, error) {
	now := time.Now().UTC()

	category := &Category{
		ID:          NewCategoryID(),
		Name:        name,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if !isActive {
		category.DeletedAt = now
	}

	return category, nil
}

func (c *Category) Update(name, description string, isActive bool) {
	if isActive {
		c.Activate()
	} else {
		c.Deactivate()
	}
	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now()
}

func (c *Category) Activate() {
	now := time.Now().UTC()

	c.DeletedAt = time.Time{}
	c.IsActive = true
	c.UpdatedAt = now
}

func (c *Category) Deactivate() {
	now := time.Now()
	if c.DeletedAt.IsZero() {
		c.DeletedAt = now
	}
	c.IsActive = false
	c.UpdatedAt = now
}

func (c *Category) Validate() error {
	var errs []error

	// Nome não pode ser nulo / vazio / whitespace
	if strings.TrimSpace(c.Name) == "" {
		errs = append(errs, errors.New(
			"category validation error: name cannot be empty or blank",
		))
	}

	// Nome mínimo 3 caracteres (após trim)
	if len(strings.TrimSpace(c.Name)) < 3 {
		errs = append(errs, errors.New(
			"category validation error: name must have at least 3 characters",
		))
	}

	if len(errs) > 0 {
		return validation.ValidationErrors{Errs: errs}
	}

	return nil
}
