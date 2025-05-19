package users_sql

const DeleteSQL = `
-- name: DeleteUser
-- Params:
--   $1: id (int64)
-- Returns: Number of rows affected
DELETE FROM users
WHERE id = $1` 