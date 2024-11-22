-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsNamesUrlsUserName :many
SELECT feeds.name, feeds.url, users.name FROM feeds
LEFT JOIN users
ON users.id = user_id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;