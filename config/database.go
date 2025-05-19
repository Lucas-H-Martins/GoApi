package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetPostgresConfig returns default PostgreSQL configuration
func GetPostgresConfig() *DBConfig {
	return &DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "goapi_db",
		SSLMode:  "disable",
	}
}

// GetConnectionString formats the connection string based on the configuration
func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// Database interface defines methods that any database implementation must satisfy
type Database interface {
	Connect() error
	Close() error
	GetDB() *sql.DB
}

// PostgresDB implements the Database interface
type PostgresDB struct {
	config *DBConfig
	db     *sql.DB
}

// NewPostgresDB creates a new PostgreSQL database instance
func NewPostgresDB(config *DBConfig) Database {
	return &PostgresDB{
		config: config,
	}
}

// Connect establishes a connection to the PostgreSQL database
func (p *PostgresDB) Connect() error {
	db, err := sql.Open("postgres", p.config.GetConnectionString())
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error pinging the database: %v", err)
	}

	p.db = db
	return nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// GetDB returns the database instance
func (p *PostgresDB) GetDB() *sql.DB {
	return p.db
}
