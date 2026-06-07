-- name: CreateUser :one
INSERT INTO users (
  username, password
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserBySessionToken :one
SELECT users.* FROM users
INNER JOIN sessions ON users.id = sessions.user_id
WHERE sessions.token = $1;

-- name: CreateSession :one
INSERT INTO sessions (
    user_id
) VALUES (
    $1
) RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions WHERE token = $1;

-- name: CreateLibrary :one
INSERT INTO libraries (
    owner_id, name, path
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetLibrariesByOwnerID :many
SELECT * FROM libraries WHERE owner_id = $1;

-- name: DeleteLibrarySecure :execrows
DELETE FROM libraries
WHERE id = $1 AND owner_id = $2;

-- name: UpdateLibrarySecure :one
UPDATE libraries
SET name = $2, path = $3
WHERE id = $1 AND owner_id = $4
RETURNING *;
