package services

import (
	"context"
	"database/sql"
	"errors"

	"goapi/models"
	"goapi/repository"
	"goapi/repository/users_sql"

	"github.com/lib/pq"
)

// UserService defines the interface for user-related business operations
type UserService interface {
	CreateUser(ctx context.Context, user *models.UserInput) (*models.UserOutput, error)
	GetUserByID(ctx context.Context, id int64) (*models.UserOutput, error)
	UpdateUser(ctx context.Context, user *models.UserOutput) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, params users_sql.SearchParams) (*models.UserList, error)
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, user *models.UserInput) (*models.UserOutput, error) {
	// Create user in repository
	createdUser, err := s.repo.Create(user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique violation
				return nil, models.ErrDuplicate
			}
		}
		return nil, err
	}

	// Return the created user's data
	return createdUser, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.UserOutput, error) {
	// get user in repository
	createdUser, err := s.repo.GetByID(int(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	// Return the created user's data
	return createdUser, nil

}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, user *models.UserOutput) error {
	// Validate user data
	if user.Name == "" {
		return models.ErrInvalidName
	}
	if user.Email == "" {
		return models.ErrInvalidEmail
	}

	// Update user in repository
	return s.repo.Update(user)
}

// DeleteUser deletes a user by their ID
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(int(id))
}

// ListUsers retrieves a list of users with pagination and filtering
func (s *userService) ListUsers(ctx context.Context, params users_sql.SearchParams) (*models.UserList, error) {
	// Validate search parameters
	if err := params.Validate(); err != nil {
		return nil, err
	}

	// Convert search params to repository params
	repoParams := repository.ListParams{
		Limit:   params.Limit,
		Offset:  params.Offset,
		Name:    params.Name,
		Email:   params.Email,
		OrderBy: params.GetOrderBy(),
	}

	// Get users from repository with pagination and filtering
	users, totalCount, err := s.repo.List(repoParams)
	if err != nil {
		return nil, err
	}

	// Convert []*models.UserOutput to []models.UserOutput
	userOutputs := make([]models.UserOutput, len(users))
	for i, user := range users {
		userOutputs[i] = *user
	}

	return &models.UserList{
		Users:      userOutputs,
		TotalCount: totalCount,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}, nil
}
