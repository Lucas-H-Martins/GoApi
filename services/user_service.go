package services

import (
	"context"

	"goapi/models"
	"goapi/repository"
	"goapi/repository/users_sql"
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
	// // Validate user data
	// if user.Name == "" {
	// 	return nil, models.ErrInvalidName
	// }
	// if user.Email == "" {
	// 	return nil, models.ErrInvalidEmail
	// }

	// Create user in repository
	createdUser, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Return the created user's data
	return createdUser, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.UserOutput, error) {
	return s.repo.GetByID(int(id))
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

	// Get users from repository
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// Apply filtering
	var filteredUsers []models.UserOutput
	for _, user := range users {
		if params.Name != "" && user.Name != params.Name {
			continue
		}
		if params.Email != "" && user.Email != params.Email {
			continue
		}
		filteredUsers = append(filteredUsers, *user)
	}

	// Apply pagination
	totalCount := int64(len(filteredUsers))
	start := params.Offset
	end := start + params.Limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}
	if start > len(filteredUsers) {
		start = len(filteredUsers)
	}

	return &models.UserList{
		Users:      filteredUsers[start:end],
		TotalCount: totalCount,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}, nil
}
