// Package http provides HTTP handlers for category operations.
package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/renamrgb/code-flix-admin-catalog/internal/application/category/create"
	"github.com/renamrgb/code-flix-admin-catalog/internal/application/category/delete"
	"github.com/renamrgb/code-flix-admin-catalog/internal/application/category/retrive"
	"github.com/renamrgb/code-flix-admin-catalog/internal/application/category/update"
)

type CategoryHandler struct {
	CreateUC *create.CreateCategoryUseCase
	UpdateUC *update.UpdateCategoryUseCase
	DeleteUC *delete.DeleteCategoryUseCase
	GetByIDUC *retrive.GetCategoryByIDUseCase
	ListUC  *retrive.ListCategoriesUseCase
}

func NewCategoryHandler(
	createUC *create.CreateCategoryUseCase,
	updateUC *update.UpdateCategoryUseCase,
	deleteUC *delete.DeleteCategoryUseCase,
	getByIDUC *retrive.GetCategoryByIDUseCase,
	listUC *retrive.ListCategoriesUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		CreateUC:  createUC,
		UpdateUC:  updateUC,
		DeleteUC:  deleteUC,
		GetByIDUC: getByIDUC,
		ListUC:    listUC,
	}
}


type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.CreateUC.Execute(create.CreateCategoryInput{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	location := fmt.Sprintf("/categories/%s", output.ID)
	w.Header().Set("Location", location)

	respondJSON(w, http.StatusCreated, output)
}


func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	output, err := h.GetByIDUC.Execute(retrive.GetCategoryByIDInput{
		ID: id,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req UpdateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	output, err := h.UpdateUC.Execute(update.UpdateCategoryInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusOK, output)
}


func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.DeleteUC.Execute(delete.DeleteCategoryInput{
		ID: id,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	input := retrive.ListCategoriesInput{
		Page:      parseInt(query.Get("page"), 1),
		PerPage:   parseInt(query.Get("per_page"), 10),
		Terms:     query.Get("terms"),
		Sort:      query.Get("sort"),
		Direction: query.Get("direction"),
	}

	output, err := h.ListUC.Execute(input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusOK, output)
}


func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}

	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return fallback
	}

	return result
}

