package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Migration struct {
	Version   string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

func runMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if _, applied := appliedMigrations[migration.Version]; !applied {
			log.Printf("Applying migration: %s", migration.Version)

			if err := applyMigration(db, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}

			log.Printf("Successfully applied migration: %s", migration.Version)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	return err
}

func loadMigrations() ([]Migration, error) {
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	migrationMap := make(map[string]*Migration)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		nameParts := strings.Split(strings.TrimSuffix(file.Name(), ".sql"), ".")
		if len(nameParts) != 2 {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join("migrations", file.Name()))
		if err != nil {
			return nil, err
		}

		if migrationMap[version] == nil {
			migrationMap[version] = &Migration{Version: version}
		}

		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationMap[version].UpSQL = string(content)
		} else if strings.HasSuffix(file.Name(), ".down.sql") {
			migrationMap[version].DownSQL = string(content)
		}
	}

	var migrations []Migration
	for _, migration := range migrationMap {
		if migration.UpSQL == "" {
			return nil, fmt.Errorf("missing up migration for version %s", migration.Version)
		}
		migrations = append(migrations, *migration)
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[string]time.Time, error) {
	rows, err := db.Query("SELECT version, applied_at FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]time.Time)
	for rows.Next() {
		var version string
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, err
		}
		applied[version] = appliedAt
	}

	return applied, nil
}

func applyMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute the migration SQL
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return err
	}

	// Record the migration as applied (this is done in the migration SQL itself)
	// But we can also do it here for consistency
	now := time.Now()
	_, err = tx.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2) ON CONFLICT (version) DO NOTHING",
		migration.Version, now)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func rollbackMigration(db *sql.DB, version string) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	var targetMigration *Migration
	for _, m := range migrations {
		if m.Version == version {
			targetMigration = &m
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %s not found", version)
	}

	if targetMigration.DownSQL == "" {
		return fmt.Errorf("no down migration available for version %s", version)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	log.Printf("Rolling back migration: %s", version)

	if _, err := tx.Exec(targetMigration.DownSQL); err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
	if err != nil {
		return err
	}

	return tx.Commit()
}
