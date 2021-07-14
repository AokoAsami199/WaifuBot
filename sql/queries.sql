-- name: getChars :many
SELECT *
FROM characters
WHERE characters.user_id = $1;
-- name: getChar :one
SELECT *
FROM characters
WHERE id = $1
    AND characters.user_id = $2;
-- name: insertChar :exec
INSERT INTO characters ("id", "user_id", "image", "name", "type")
VALUES ($1, $2, $3, $4, $5);
-- name: giveChar :one
UPDATE characters
SET "type" = 'TRADE',
    "user_id" = @given
WHERE characters.id = @id
    AND characters.user_id = @giver
RETURNING *;
-- name: createUser :exec
INSERT INTO users (user_id)
VALUES ($1);