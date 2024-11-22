-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
        )
    RETURNING *
)
SELECT
    inserted_feed_follows.*, 
    users.name AS user_name, 
    feeds.name AS feed_name 
FROM inserted_feed_follows
INNER JOIN users
ON users.id = user_id
INNER JOIN feeds
ON feeds.id = feed_id
;

-- name: GetFeedFollowForUser :many
SELECT 
    feed_follows.*,
    users.name AS user_name,
    feeds.name AS feed_name
 FROM feed_follows
 INNER JOIN users 
    ON users.id = feed_follows.user_id
INNER JOIN feeds 
    ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
;