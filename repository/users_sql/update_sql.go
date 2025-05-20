package users_sql

const UpdateSQL = `
-- name: UpdateUser
-- Params:
--   $1: name (string)
--   $2: email (string)
--   $3: updated_at (timestamp)
--   $4: id (int64)
-- Returns: Number of rows affected
UPDATE users
SET
    name = $1,
    email = $2,
    updated_at = $3
WHERE 
    id = $4 AND 
    deleted_at IS NULL`
