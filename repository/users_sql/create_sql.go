package users_sql

const CreateSQL = `
-- name: CreateUser
-- Params:
--   $1: name (string)
--   $2: email (string)
--   $3: created_at (timestamp)
--   Note: created_at is also used for updated_at initially
INSERT INTO users (
    name,
    email,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $3
)
RETURNING id` 