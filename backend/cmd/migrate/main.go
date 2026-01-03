package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"github.com/portfolio/backend/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	var direction = flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrations table
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	migrationsDir := "./migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), fmt.Sprintf(".%s.sql", *direction)) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, filename := range migrationFiles {
		version := extractVersion(filename)
		if *direction == "up" {
			if alreadyApplied(db, version) {
				log.Printf("Migration %s already applied, skipping", filename)
				continue
			}
			if err := runMigration(db, filepath.Join(migrationsDir, filename), version); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			log.Printf("Applied migration: %s", filename)
		}
	}

	log.Println("Migrations completed successfully!")
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func extractVersion(filename string) string {
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		return parts[0]
	}
	return filename
}

func alreadyApplied(db *sql.DB, version string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func runMigration(db *sql.DB, filepath, version string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(sqlBytes)); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}

	return tx.Commit()
}

