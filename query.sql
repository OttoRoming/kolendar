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

-- name: DeleteCalendarByIDAndUserID :execrows
DELETE FROM calendars WHERE id = $1 AND user_id = $2;

-- name: UpdateCalendarByIDAndUserID :one
UPDATE calendars SET name = $3, color = $4 WHERE id = $1 AND user_id = $2 RETURNING *;

-- name: GetEventsByCalendarIDAndUserID :many
SELECT events.* FROM events
INNER JOIN calendars ON events.calendar_id = calendars.id
WHERE calendars.user_id = $2 
  AND events.calendar_id = $1;

-- name: CreateEventSecure :one
INSERT INTO events(
    calendar_id, title, description, all_day, start_time, end_time, location
) SELECT $1, $2, $3, $4, $5, $6, $7
FROM calendars
WHERE calendars.id = $1 AND calendars.user_id = $8
RETURNING *;

-- name: DeleteEventByIDAndUserID :execrows
DELETE FROM events
USING calendars
WHERE events.id = $1 
  AND events.calendar_id = calendars.id 
  AND calendars.user_id = $2;

-- name: UpdateEventByIDAndUserID :one
UPDATE events
SET 
    calendar_id = $1,
    title = $2,
    description = $3,
    all_day = $4,
    start_time = $5,
    end_time = $6,
    location = $7
FROM
    calendars AS old_calendar,
    calendars AS new_calendar
WHERE events.id = $8
  -- 1. make sure the user owns the OLD calendar that the event is currently in
  AND events.calendar_id = old_calendar.id
  AND old_calendar.user_id = $9
  -- 2. make sure the user owns the NEW calendar that the event will be moved to
  AND new_calendar.id = $1
  AND new_calendar.user_id = $9
RETURNING events.*;

