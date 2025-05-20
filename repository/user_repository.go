package repository

import (
	"database/sql"
	"goapi/logger"
	"goapi/models"
	"goapi/repository/users_sql"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id int) (*models.User, error)
	List() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

// PostgresUserRepository implements UserRepository for PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create implements the Create method of UserRepository
func (r *PostgresUserRepository) Create(user *models.User) (*models.User, error) {
	query := users_sql.CreateUserSQL
	var userResponse models.User
	// Execute the query and scan the result into the userResponse struct
	err := r.db.QueryRow(query, user.Name, user.Email, user.CreatedAt).Scan(&userResponse.ID, &userResponse.Name, &userResponse.Email, &userResponse.CreatedAt, &userResponse.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &userResponse, nil
}

// GetByID implements the GetByID method of UserRepository
func (r *PostgresUserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email FROM users WHERE id = $1`

	logger.Debug("Executing query: %s with id: %d", query, id)

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		logger.Error("Error retrieving user with id %d: %v", id, err)
		return nil, err
	}
	logger.Debug("User retrieved: %+v", user)
	return user, nil
}

// List implements the List method of UserRepository
func (r *PostgresUserRepository) List() ([]*models.User, error) {
	query := `SELECT id, name, email FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Update implements the Update method of UserRepository
func (r *PostgresUserRepository) Update(user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	return err
}

// Delete implements the Delete method of UserRepository
func (r *PostgresUserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
