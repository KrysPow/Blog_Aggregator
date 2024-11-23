-- name: CreatePost :exec
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7, 
    $8
);

-- name: GetPostsForUser :many
Select posts.* 
FROM posts
INNER JOIN feed_follows
    ON posts.feed_id = feed_follows.feed_id
INNER JOIN users
    ON feed_follows.user_id = users.id
WHERE users.name = $1
ORDER BY posts.updated_at ASC 
LIMIT $2;