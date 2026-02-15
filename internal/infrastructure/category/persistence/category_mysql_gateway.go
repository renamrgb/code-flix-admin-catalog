// Package persistence provides MySQL gateway implementations for category data access.
package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/pagination"
)

type MySQLCategoryGateway struct {
	DB *sql.DB
}

func NewMySQLCategoryGateway(db *sql.DB) *MySQLCategoryGateway {
	return &MySQLCategoryGateway{DB: db}
}

func (g *MySQLCategoryGateway) CreateCategory(cat *category.Category) (*category.Category, error) {
	query := `
		INSERT INTO categories (id, name, description, activated, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := g.DB.Exec(
		query,
		cat.ID.String(),
		cat.Name,
		cat.Description,
		cat.IsActive,
		cat.CreatedAt,
		cat.UpdatedAt,
		nullTime(cat.DeletedAt),
	)

	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (g *MySQLCategoryGateway) GetCategoryByID(id category.CategoryID) (*category.Category, error) {
	query := `
		SELECT id, name, description, activated, created_at, updated_at, deleted_at
		FROM categories
		WHERE id = ?
	`

	row := g.DB.QueryRow(query, id.String())

	var cat category.Category
	var deletedAt sql.NullTime

	err := row.Scan(
		&cat.ID,
		&cat.Name,
		&cat.Description,
		&cat.IsActive,
		&cat.CreatedAt,
		&cat.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	if deletedAt.Valid {
		cat.DeletedAt = deletedAt.Time
	}

	return &cat, nil
}

func (g *MySQLCategoryGateway) UpdateCategory(cat *category.Category) (*category.Category, error) {
	query := `
		UPDATE categories
		SET name = ?, description = ?, activated = ?, updated_at = ?, deleted_at = ?
		WHERE id = ?
	`

	_, err := g.DB.Exec(
		query,
		cat.Name,
		cat.Description,
		cat.IsActive,
		cat.UpdatedAt,
		nullTime(cat.DeletedAt),
		cat.ID.String(),
	)

	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (g *MySQLCategoryGateway) DeleteCategory(id category.CategoryID) error {
	query := `DELETE FROM categories WHERE id = ?`
	_, err := g.DB.Exec(query, id.String())
	return err
}

func (g *MySQLCategoryGateway) FindAll(query category.SearchCategoryQuery) (*pagination.Pagination[category.Category], error) {
	offset := (query.Page - 1) * query.PerPage

	whereClause := ""
	args := []any{}

	if query.Terms != "" {
		whereClause = "WHERE name LIKE ? OR description LIKE ?"
		terms := "%" + query.Terms + "%"
		args = append(args, terms, terms)
	}

	sort := resolveSort(query.Sort)
	direction := resolveDirection(query.Direction)

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM categories %s`, whereClause)

	var total int
	if err := g.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	searchQuery := fmt.Sprintf(`
		SELECT id, name, description, activated, created_at, updated_at, deleted_at
		FROM categories
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereClause, sort, direction)

	args = append(args, query.PerPage, offset)

	rows, err := g.DB.Query(searchQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []category.Category

	for rows.Next() {
		var cat category.Category
		var deletedAt sql.NullTime

		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&cat.IsActive,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}

		if deletedAt.Valid {
			cat.DeletedAt = deletedAt.Time
		}

		categories = append(categories, cat)
	}

	return &pagination.Pagination[category.Category]{
		CurrentPage: query.Page,
		PerPage:     query.PerPage,
		Total:       total,
		Items:       categories,
	}, nil
}

func resolveSort(sort string) string {
	switch sort {
	case "name", "created_at", "updated_at":
		return sort
	default:
		return "created_at"
	}
}

func resolveDirection(direction string) string {
	if strings.ToLower(direction) == "desc" {
		return "DESC"
	}
	return "ASC"
}

func nullTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}
