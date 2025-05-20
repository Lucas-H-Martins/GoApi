package users_sql

import (
	"fmt"
)

const ListSQL = `
-- name: ListUsers
-- Params:
--   $1: limit (int)
--   $2: offset (int)
-- Optional Params (used in WHERE clause):
--   $3: name (string) - for ILIKE search
--   $4: email (string) - for ILIKE search
-- Returns: Multiple rows of user data with pagination
WITH filtered_users AS (
    SELECT
        id,
        name,
        email,
        created_at,
        updated_at
    FROM users
    WHERE
        deleted_at IS NULL
        AND (
            ($3 IS NULL OR $3 = '') AND ($4 IS NULL OR $4 = '')
            OR ($3 IS NOT NULL AND $3 != '' AND name ILIKE '%' || $3 || '%')
            OR ($4 IS NOT NULL AND $4 != '' AND email ILIKE '%' || $4 || '%')
        )
)
SELECT
    id,
    name,
    email,
    created_at,
    updated_at,
    (SELECT COUNT(*) FROM filtered_users) as total_count
FROM filtered_users
%s -- This will be replaced with ORDER BY clause
LIMIT $1
OFFSET $2`

// GetListSQL returns the formatted SQL with the ORDER BY clause
func GetListSQL(orderBy string) string {
	return fmt.Sprintf(ListSQL, "ORDER BY "+orderBy)
}

const CountSQL = `
-- name: CountUsers
-- Params:
--   $1: name (string) - for ILIKE search
--   $2: email (string) - for ILIKE search
-- Returns: Total count of matching users
SELECT COUNT(*)
FROM users
WHERE
    deleted_at IS NULL
    AND CASE 
        WHEN $1 != '' THEN name ILIKE '%' || $1 || '%'
        WHEN $2 != '' THEN email ILIKE '%' || $2 || '%'
        ELSE true
    END`
