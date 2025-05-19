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
	files                     []MigrationFile
	db                        *sql.DB
	defaultUpMigrationsPath   string
	defaultDownMigrationsPath string
}

func NewManager(db *sql.DB) (*Manager, error) {
	m := &Manager{db: db}

	m.defaultUpMigrationsPath = "./migrations/up"
	m.defaultDownMigrationsPath = "./migrations/down"

	if err := m.createMigrationsTable(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			undone_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
		)
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Manager) isMigrationUpExecuted(name string) (bool, error) {
	var exists bool
	err := m.db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = $1 and undone_at is null)", name).Scan(&exists)
	return exists, err
}

func (m *Manager) RunMigrationsUp() error {
	for _, file := range m.files {

		file.PrintMigrations("up")
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
	for _, file := range m.files {

		file.PrintMigrations("Down")
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
		_, err = tx.Exec("UPDATE migrations SET undone_at = NOW() WHERE name = $1 AND undone_at IS NULL", file.Name)

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

func (m *Manager) LoadMigrationsUp(dir string) error {

	// Read all files in the migrations directory
	files, err := os.ReadDir(m.defaultUpMigrationsPath)

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

		executed, err := m.isMigrationUpExecuted(file.Name())
		if err != nil {
			return fmt.Errorf("error checking migration status: %w", err)
		}
		if executed {
			logger.Debug("Migration %s already executed, skipping...", file.Name())
			continue
		}

		// Read file content
		content, err := os.ReadFile(filepath.Join(m.defaultUpMigrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
		}

		migrationFile := MigrationFile{
			Path:    filepath.Join(m.defaultUpMigrationsPath, file.Name()),
			Name:    file.Name(),
			Content: string(content),
		}
		m.files = append(m.files, migrationFile)
	}

	return nil
}

func (m *Manager) LoadMigrationsDown(dir string) error {

	// Read all files in the migrations directory
	files, err := os.ReadDir(m.defaultDownMigrationsPath)

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

		exist, err := m.isMigrationUpExecuted(file.Name())

		if err != nil {
			return fmt.Errorf("error checking migration status: %w", err)
		}
		if !exist {
			logger.Debug("Migration down %s already executed, skipping...", file.Name())
			continue
		}

		// Read file content
		content, err := os.ReadFile(filepath.Join(m.defaultDownMigrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
		}

		migrationFile := MigrationFile{
			Path:    filepath.Join(m.defaultDownMigrationsPath, file.Name()),
			Name:    file.Name(),
			Content: string(content),
		}
		m.files = append(m.files, migrationFile)
	}

	return nil
}

func (file *MigrationFile) PrintMigrations(dir string) {
	if dir == "down" {
		logger.Debug("=== Down Migration ===")
	} else if dir == "up" {
		logger.Debug("=== Up Migration ===")
	}
	logger.Debug("File: %s", file.Name)
	logger.Debug("----------------------------------------")
}
