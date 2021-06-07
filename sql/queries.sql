-- name: GetChars :many
SELECT *
FROM characters
WHERE characters.user_id = $1;

-- name: GetChar :one
SELECT *
FROM characters
WHERE id = $1
    AND characters.user_id = $2;

-- name: InsertChar :exec
INSERT INTO characters ("id", "user_id", "image", "name", "type")
VALUES ($1, $2, $3, $4, $5);

-- name: GiveChar :exec
UPDATE characters
SET user_id = $3
WHERE characters.id = $1
    AND characters.user_id = $2;

-- name: CreateUser :exec
INSERT INTO users (user_id)
VALUES ($1);