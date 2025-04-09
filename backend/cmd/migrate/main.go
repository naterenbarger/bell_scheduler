package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"bell_scheduler/internal/config"
	"bell_scheduler/internal/models"

	"gorm.io/gorm"
)

var (
	dbPath = flag.String("db", "data/bell_scheduler.db", "Path to the database file")
)

func main() {
	flag.Parse()

	// Create migrations directory if it doesn't exist
	migrationsDir := "internal/store/migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		fmt.Printf("Failed to create migrations directory: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := config.NewDB(*dbPath)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Get command from arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [command]")
		fmt.Println("Commands:")
		fmt.Println("  create [name]  - Create a new migration")
		fmt.Println("  up            - Run all pending migrations")
		fmt.Println("  down          - Rollback the last migration")
		fmt.Println("  status        - Show migration status")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a migration name")
			os.Exit(1)
		}
		createMigration(db, os.Args[2])
	case "up":
		runMigrations(db, true)
	case "down":
		runMigrations(db, false)
	case "status":
		showStatus(db)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func createMigration(db *gorm.DB, name string) {
	// Create migration file
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, strings.ToLower(name))
	filepath := filepath.Join("internal/store/migrations", filename)

	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Failed to create migration file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write migration template
	template := `-- Migration: %s
-- Created at: %s

-- Up Migration
BEGIN;

-- Add your migration SQL here

COMMIT;

-- Down Migration
BEGIN;

-- Add your rollback SQL here

COMMIT;
`
	fmt.Fprintf(file, template, name, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Created migration: %s\n", filename)
}

func runMigrations(db *gorm.DB, up bool) {
	// Get list of migration files
	files, err := filepath.Glob("internal/store/migrations/*.sql")
	if err != nil {
		fmt.Printf("Failed to list migration files: %v\n", err)
		os.Exit(1)
	}

	// Sort files by name (timestamp)
	sort.Strings(files)

	// Create migrations table if it doesn't exist
	err = db.AutoMigrate(&models.Migration{})
	if err != nil {
		fmt.Printf("Failed to create migrations table: %v\n", err)
		os.Exit(1)
	}

	if up {
		// Run pending migrations
		for _, file := range files {
			var migration models.Migration
			filename := filepath.Base(file)
			if err := db.Where("filename = ?", filename).First(&migration).Error; err == nil {
				continue // Migration already applied
			}

			// Read migration file
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("Failed to read migration file %s: %v\n", file, err)
				continue
			}

			// Split into up and down migrations
			parts := strings.Split(string(content), "-- Down Migration")
			if len(parts) != 2 {
				fmt.Printf("Invalid migration file format: %s\n", file)
				continue
			}

			// Extract SQL statements
			upSQL := extractSQL(parts[0])
			_ = extractSQL(parts[1]) // Store down migration for later use

			// Run migration
			err = db.Transaction(func(tx *gorm.DB) error {
				// Execute up migration
				for _, sql := range upSQL {
					if err := tx.Exec(sql).Error; err != nil {
						return fmt.Errorf("failed to execute migration %s: %v", file, err)
					}
				}

				// Record migration
				return tx.Create(&models.Migration{
					Filename:  filename,
					AppliedAt: time.Now(),
				}).Error
			})

			if err != nil {
				fmt.Printf("Failed to apply migration %s: %v\n", file, err)
				os.Exit(1)
			}

			fmt.Printf("Applied migration: %s\n", filename)
		}
	} else {
		// Rollback last migration
		var migration models.Migration
		if err := db.Order("applied_at desc").First(&migration).Error; err != nil {
			fmt.Println("No migrations to rollback")
			return
		}

		// Read migration file
		file := filepath.Join("internal/store/migrations", migration.Filename)
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read migration file %s: %v\n", file, err)
			os.Exit(1)
		}

		// Split into up and down migrations
		parts := strings.Split(string(content), "-- Down Migration")
		if len(parts) != 2 {
			fmt.Printf("Invalid migration file format: %s\n", file)
			os.Exit(1)
		}

		// Extract SQL statements
		downSQL := extractSQL(parts[1])

		// Run rollback
		err = db.Transaction(func(tx *gorm.DB) error {
			// Execute down migration
			for _, sql := range downSQL {
				if err := tx.Exec(sql).Error; err != nil {
					return fmt.Errorf("failed to execute rollback %s: %v", file, err)
				}
			}

			// Remove migration record
			return tx.Delete(&migration).Error
		})

		if err != nil {
			fmt.Printf("Failed to rollback migration %s: %v\n", file, err)
			os.Exit(1)
		}

		fmt.Printf("Rolled back migration: %s\n", migration.Filename)
	}
}

func showStatus(db *gorm.DB) {
	// Get list of migration files
	files, err := filepath.Glob("internal/store/migrations/*.sql")
	if err != nil {
		fmt.Printf("Failed to list migration files: %v\n", err)
		os.Exit(1)
	}

	// Sort files by name (timestamp)
	sort.Strings(files)

	// Get applied migrations
	var appliedMigrations []models.Migration
	if err := db.Find(&appliedMigrations).Error; err != nil {
		fmt.Printf("Failed to get applied migrations: %v\n", err)
		os.Exit(1)
	}

	// Create map of applied migrations
	applied := make(map[string]bool)
	for _, m := range appliedMigrations {
		applied[m.Filename] = true
	}

	// Show status
	fmt.Println("Migration Status:")
	fmt.Println("----------------")
	for _, file := range files {
		filename := filepath.Base(file)
		status := "Pending"
		if applied[filename] {
			status = "Applied"
		}
		fmt.Printf("%s: %s\n", filename, status)
	}
}

func extractSQL(content string) []string {
	// Split by semicolon and clean up
	statements := strings.Split(content, ";")
	var sql []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" && !strings.HasPrefix(stmt, "--") {
			sql = append(sql, stmt)
		}
	}
	return sql
}
