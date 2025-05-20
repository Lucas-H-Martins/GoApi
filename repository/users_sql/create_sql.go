package users_sql

const CreateUserSQL = `
-- name: CreateUser
-- Params:
--   $1: name (string)
--   $2: email (string)
INSERT INTO users (
    name,
    email,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    now(),
    now()
)
RETURNING 
    id,
    name,
    email,
    created_at,
    updated_at
    `
