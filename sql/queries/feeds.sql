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

-- name: GetFeeds :many
select feeds.name, feeds.url, users.name
from feeds
inner join users
  on feeds.user_id = users.id;

-- name: GetFeedByURL :one
select * from feeds
where url = $1;

-- name: MarkFeedFetched :one
update feeds
set last_fetched_at = NOW(), updated_at = NOW()
where id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at nulls first
limit 1;

