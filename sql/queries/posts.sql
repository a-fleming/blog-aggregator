-- name: CreatePost :one
INSERT INTO posts (
    title,
    url,
    description,
    published_at,
    feed_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
	posts.title AS title,
	posts.url AS url,
	posts.description AS description,
	posts.id AS post_id,
	posts.created_at AS created_at,
	posts.updated_at AS updated_at,
	posts.published_at AS published_at,
	posts.feed_id AS feed_id
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
JOIN users ON feed_follows.user_id = users.id
WHERE users.id = $1
ORDER BY published_at DESC
LIMIT $2;
