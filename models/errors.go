package models

import "errors"

var (
	// ErrInvalidName is returned when a user's name is invalid
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidEmail is returned when a user's email is invalid
	ErrInvalidEmail = errors.New("invalid email")
)
