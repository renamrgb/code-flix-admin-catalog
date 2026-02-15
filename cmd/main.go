package main

import (
	"log"
	"net/http"

	createUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/create"
	deleteUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/delete"
	retriveUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/retrive"
	updateUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/update"

	categoryHTTP "github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/interfaces/http"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/category/persistence"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/database/migration"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/database/mysql"
	// Update the import path below to match the actual location of your category handler package.
)

func main() {
	cfg, err := mysql.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, err := mysql.NewConnection(cfg)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	log.Println("Database connected")

	if err := migration.RunMigrations(db, cfg.Database); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	log.Println("Migrations executed successfully")

	var gateway category.CategoryGateway = persistence.NewMySQLCategoryGateway(db)

	createUseCase := createUC.NewCreateCategoryUseCase(gateway)
	updateUseCase := updateUC.NewUpdateCategoryUseCase(gateway)
	deleteUseCase := deleteUC.NewDeleteCategoryUseCase(gateway)
	getByIDUseCase := retriveUC.NewGetCategoryByIDUseCase(gateway)
	listUseCase := retriveUC.NewListCategoriesUseCase(gateway)

	handler := categoryHTTP.NewCategoryHandler(
		createUseCase,
		updateUseCase,
		deleteUseCase,
		getByIDUseCase,
		listUseCase,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /categories", handler.CreateCategory)
	mux.HandleFunc("GET /categories", handler.ListCategories)
	mux.HandleFunc("GET /categories/{id}", handler.GetCategoryByID)
	mux.HandleFunc("PUT /categories/{id}", handler.UpdateCategory)
	mux.HandleFunc("DELETE /categories/{id}", handler.DeleteCategory)

	log.Println("HTTP server running at :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
