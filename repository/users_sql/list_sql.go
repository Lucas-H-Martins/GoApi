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
        CASE 
            WHEN $3 != '' THEN name ILIKE '%' || $3 || '%'
            WHEN $4 != '' THEN email ILIKE '%' || $4 || '%'
            ELSE true
        END
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
    CASE 
        WHEN $1 != '' THEN name ILIKE '%' || $1 || '%'
        WHEN $2 != '' THEN email ILIKE '%' || $2 || '%'
        ELSE true
    END` 