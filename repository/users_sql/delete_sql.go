package users_sql

const DeleteSQL = `
-- name: SoftDeleteUser
-- Params:
--   $1: id (int64)
-- Returns: Number of rows affected
UPDATE users
SET deleted_at = NOW()
WHERE id = $1`
