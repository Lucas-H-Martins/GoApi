package repository

import (
	"database/sql"
	"goapi/logger"
	"goapi/models"
	"goapi/repository/users_sql"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.UserInput) (*models.UserOutput, error)
	GetByID(id int) (*models.UserOutput, error)
	List(params ListParams) ([]*models.UserOutput, int64, error)
	Update(user *models.UserOutput) error
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
func (r *PostgresUserRepository) Create(user *models.UserInput) (*models.UserOutput, error) {
	query := users_sql.CreateUserSQL
	var userResponse models.UserOutput
	// Execute the query and scan the result into the userResponse struct
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(&userResponse.ID, &userResponse.Name, &userResponse.Email, &userResponse.CreatedAt, &userResponse.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &userResponse, nil
}

// GetByID implements the GetByID method of UserRepository
func (r *PostgresUserRepository) GetByID(id int) (*models.UserOutput, error) {
	user := &models.UserOutput{}
	query := users_sql.GetByIDSQL
	logger.Debug("Executing query: %s with id: %d", query, id)

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error("Error retrieving user with id %d: %v", id, err)
		return nil, err
	}
	logger.Debug("User retrieved: %+v", user)
	return user, nil
}

// ListParams represents the parameters for listing users
type ListParams struct {
	Limit   int
	Offset  int
	Name    string
	Email   string
	OrderBy string
}

// List implements the List method of UserRepository
func (r *PostgresUserRepository) List(params ListParams) ([]*models.UserOutput, int64, error) {
	query := users_sql.GetListSQL(params.OrderBy)
	rows, err := r.db.Query(query, params.Limit, params.Offset, params.Name, params.Email)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.UserOutput
	var totalCount int64

	for rows.Next() {
		user := &models.UserOutput{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&totalCount,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// Update implements the Update method of UserRepository
func (r *PostgresUserRepository) Update(user *models.UserOutput) error {
	query := users_sql.UpdateSQL
	_, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	return err
}

// Delete implements the Delete method of UserRepository
func (r *PostgresUserRepository) Delete(id int) error {
	query := users_sql.DeleteSQL
	_, err := r.db.Exec(query, id)
	return err
}
