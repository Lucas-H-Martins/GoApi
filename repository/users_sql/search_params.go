package users_sql

import "fmt"

type SortOrder string

const (
	ASC  SortOrder = "ASC"
	DESC SortOrder = "DESC"
)

type SearchParams struct {
	Limit  int
	Offset int
	Order  SortOrder
	Name   string
	Email  string
}

const (
	MinLimit = 1
	MaxLimit = 50
)

// Validate checks if the search parameters are valid
func (p *SearchParams) Validate() error {
	if p.Limit < MinLimit || p.Limit > MaxLimit {
		return fmt.Errorf("limit must be between %d and %d", MinLimit, MaxLimit)
	}

	if p.Offset < 0 {
		return fmt.Errorf("offset must be greater than or equal to 0")
	}

	if p.Order != "" && p.Order != ASC && p.Order != DESC {
		return fmt.Errorf("order must be either ASC or DESC")
	}

	if p.Name != "" && p.Email != "" {
		return fmt.Errorf("cannot search by both name and email")
	}

	return nil
}

// GetOrderBy returns the ORDER BY clause based on the search parameters
func (p *SearchParams) GetOrderBy() string {
	if p.Order == "" {
		p.Order = ASC
	}

	// Default to ordering by id if no specific field is being searched
	orderField := "id"
	if p.Name != "" {
		orderField = "name"
	} else if p.Email != "" {
		orderField = "email"
	}

	return fmt.Sprintf("%s %s", orderField, p.Order)
}
