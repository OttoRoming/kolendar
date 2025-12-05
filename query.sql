-- name: CreateUser :one
INSERT INTO users (
  id, username, password
) VALUES (
  ?, ?, ?
) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: CreateCalendar :one
INSERT INTO calendars (
  id, user_id, name, color
) VALUES (
    ?, ?, ?, ?
) RETURNING *;

-- name: GetCalendarsByUserID :many
SELECT * FROM calendars WHERE user_id = ?;
