package repository

import (
	"context"
	"database/sql"
	"goapi/logger"
	"goapi/models"
	"goapi/repository/users_sql"
	"time"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params users_sql.SearchParams) (*models.UserList, error)
}

// PostgresUserRepository implements UserRepository for PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	logger.Debug("Initializing PostgreSQL user repository")
	return &PostgresUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	logger.Debug("Creating new user with name: %s, email: %s", user.Name, user.Email)
	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		users_sql.CreateSQL,
		user.Name,
		user.Email,
		now,
	).Scan(&user.ID)

	if err != nil {
		logger.Error("Failed to create user: %v", err)
		return err
	}

	logger.Info("Created user with ID: %d", user.ID)
	return nil
}

// GetByID retrieves a user by their ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	logger.Debug("Getting user by ID: %d", id)
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, users_sql.GetByIDSQL, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		logger.Debug("No user found with ID: %d", id)
		return nil, nil
	}
	if err != nil {
		logger.Error("Error getting user by ID %d: %v", id, err)
		return nil, err
	}

	logger.Debug("Found user: %+v", user)
	return user, nil
}

// Update modifies an existing user in the database
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	logger.Debug("Updating user ID: %d with name: %s, email: %s", user.ID, user.Name, user.Email)
	result, err := r.db.ExecContext(
		ctx,
		users_sql.UpdateSQL,
		user.Name,
		user.Email,
		time.Now(),
		user.ID,
	)
	if err != nil {
		logger.Error("Error updating user %d: %v", user.ID, err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected for user update %d: %v", user.ID, err)
		return err
	}
	if rows == 0 {
		logger.Warn("No user found to update with ID: %d", user.ID)
		return sql.ErrNoRows
	}

	logger.Info("Successfully updated user ID: %d", user.ID)
	return nil
}

// Delete removes a user from the database
func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	logger.Debug("Deleting user with ID: %d", id)
	result, err := r.db.ExecContext(ctx, users_sql.DeleteSQL, id)
	if err != nil {
		logger.Error("Error deleting user %d: %v", id, err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected for user deletion %d: %v", id, err)
		return err
	}
	if rows == 0 {
		logger.Warn("No user found to delete with ID: %d", id)
		return sql.ErrNoRows
	}

	logger.Info("Successfully deleted user ID: %d", id)
	return nil
}

// List returns paginated users from the database with optional filtering
func (r *PostgresUserRepository) List(ctx context.Context, params users_sql.SearchParams) (*models.UserList, error) {
	logger.Debug("Listing users with params: %+v", params)
	
	// Validate search parameters
	if err := params.Validate(); err != nil {
		logger.Error("Invalid search parameters: %v", err)
		return nil, err
	}

	// Get the formatted SQL with ORDER BY clause
	query := users_sql.GetListSQL(params.GetOrderBy())

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query,
		params.Limit,
		params.Offset,
		params.Name,
		params.Email,
	)
	if err != nil {
		logger.Error("Error executing list query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	var totalCount int64

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&totalCount,
		)
		if err != nil {
			logger.Error("Error scanning user row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error iterating through rows: %v", err)
		return nil, err
	}

	logger.Debug("Found %d users (total count: %d)", len(users), totalCount)
	return &models.UserList{
		Users:      users,
		TotalCount: totalCount,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}, nil
} 