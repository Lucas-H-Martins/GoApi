package models

type UserInput struct {
	Name  string `json:"name" binding:"required,min=3"`
	Email string `json:"email" binding:"required,email"`
}

type UserOutput struct {
	ID        *int64 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserList struct {
	Users      []UserOutput `json:"users"`
	TotalCount int64        `json:"total_count"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
}
