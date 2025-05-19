package migrations

import (
	"database/sql"
	"fmt"
	"goapi/logger"
	"os"
	"path/filepath"
)

type MigrationFile struct {
	Path    string
	Name    string
	Content string
}

type Manager struct {
	files []MigrationFile
	db    *sql.DB
}

func NewManager(db *sql.DB) (*Manager, error) {
	m := &Manager{db: db}

	if err := m.createMigrationsTable(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Manager) isMigrationExecuted(name string) (bool, error) {
	var exists bool
	err := m.db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = $1)", name).Scan(&exists)
	return exists, err
}

func (m *Manager) RunMigrationsUp() error {
	for _, file := range m.files {

		file.PrintUpMigrations()
		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		// Execute migration
		_, err = tx.Exec(file.Content)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %w", file.Name, err)
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", file.Name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %w", file.Name, err)
		}

		// Commit transaction
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %s: %w", file.Name, err)
		}

		logger.Info("Successfully executed migration: %s", file.Name)
	}

	return nil
}

func (m *Manager) RunMigrationsDown() error {
	// future implement
	return nil
}

func (m *Manager) LoadMigrations() error {

	// Load environment file from main folder
	const defaultMigrationsPath = "./migrations/up"
	// Read all files in the migrations directory
	files, err := os.ReadDir(defaultMigrationsPath)

	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Process each file
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		executed, err := m.isMigrationExecuted(file.Name())
		if err != nil {
			return fmt.Errorf("error checking migration status: %w", err)
		}
		if executed {
			logger.Debug("Migration %s already executed, skipping...", file.Name())
			continue
		}

		// Read file content
		content, err := os.ReadFile(filepath.Join(defaultMigrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
		}

		migrationFile := MigrationFile{
			Path:    filepath.Join(defaultMigrationsPath, file.Name()),
			Name:    file.Name(),
			Content: string(content),
		}
		m.files = append(m.files, migrationFile)
	}

	return nil
}

func (file *MigrationFile) PrintUpMigrations() {
	logger.Debug("=== Up Migration ===")
	logger.Debug("File: %s", file.Name)
	logger.Debug("----------------------------------------")
}
