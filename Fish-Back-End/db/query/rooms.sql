-- name: ListRooms :many
SELECT * FROM rooms ORDER BY min_bet ASC;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;