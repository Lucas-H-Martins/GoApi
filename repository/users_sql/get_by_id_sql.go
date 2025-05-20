package users_sql

const GetByIDSQL = `
-- name: GetUserByID
-- Params:
--   $1: id (int64)
-- Returns: Single row with user data
SELECT
    id,
    name,
    email,
    created_at,
    updated_at
FROM users
WHERE 
    id = $1 AND 
    deleted_at IS NULL`
