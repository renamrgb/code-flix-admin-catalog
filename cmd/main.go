package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	createCategoryUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/create"
	deleteCategoryUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/delete"
	retriveCategoryUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/retrive"
	updateCategoryUC "github.com/renamrgb/code-flix-admin-catalog/internal/application/category/update"

	categoryHTTP "github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/interfaces/http"

	"github.com/renamrgb/code-flix-admin-catalog/internal/domain/category"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/category/persistence"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/database/migration"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/database/mysql"
	// Update the import path below to match the actual location of your category handler package.
)

func main() {
	_ = godotenv.Load()
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

	createUseCase := createCategoryUC.NewCreateCategoryUseCase(gateway)
	updateUseCase := updateCategoryUC.NewUpdateCategoryUseCase(gateway)
	deleteUseCase := deleteCategoryUC.NewDeleteCategoryUseCase(gateway)
	getByIDUseCase := retriveCategoryUC.NewGetCategoryByIDUseCase(gateway)
	listUseCase := retriveCategoryUC.NewListCategoriesUseCase(gateway)

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
