package models

import "time"

type User struct {
	ID        *int64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserList struct {
	Users      []User `json:"users"`
	TotalCount int64  `json:"total_count"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}
