-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
WITH users_feeds AS(
    SELECT feed_follows.*, feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
INNER JOIN users
    ON feed_follows.user_id = users.id
WHERE users.name = $1
)
SELECT posts.*
FROM posts
INNER JOIN users_feeds
    ON posts.feed_id = users_feeds.feed_id
ORDER BY posts.published_at DESC
LIMIT $2;