-- name: CreateFeedFollow :one

WITH new_feed_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)

SELECT new_feed_follow.*, 
feeds.name AS feed_name,
users.name AS user_name
FROM new_feed_follow
INNER JOIN feeds ON feeds.id = new_feed_follow.feed_id
INNER JOIN users ON users.id = new_feed_follow.user_id;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id=$2;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, 
       feeds.name as feed_name,
       feeds.url as feed_url
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;