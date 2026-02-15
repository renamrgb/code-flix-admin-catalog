package main

import (
	"log"

	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/database/migration"
	"github.com/renamrgb/code-flix-admin-catalog/internal/infrastructure/database/mysql"
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
}
