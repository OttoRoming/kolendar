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

-- name: CreateSession :one
INSERT INTO sessions (
    user_id
) VALUES (
    $1
) RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions WHERE token = $1;

-- name: CreateCalendar :one
INSERT INTO calendars (
  user_id, name, color
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCalendarsByUserID :many
SELECT * FROM calendars WHERE user_id = $1;

-- name: GetCalendarByID :one
SELECT * FROM calendars WHERE id = $1;

-- name: DeleteCalendarByID :exec
DELETE FROM calendars WHERE id = $1;

-- name: UpdateCalendarByID :one
UPDATE calendars SET name = $2, color = $3 WHERE id = $1 RETURNING *;

-- name: GetEventsByCalendarID :many
SELECT * FROM events WHERE calendar_id = $1;

-- name: CreateEvent :one
INSERT INTO events (
    calendar_id, title, description, all_day, start_time, end_time, location
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

