-- name: ListFishes :many
SELECT * FROM fishes ORDER BY reward_multiplier ASC;

-- name: GetFishByID :one
SELECT * FROM fishes WHERE id = $1;