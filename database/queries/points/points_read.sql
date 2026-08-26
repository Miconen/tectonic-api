-- name: GetGuildPointSources :many
SELECT "source", "points", "name" FROM point_sources WHERE guild_id = @guild_id;

-- name: GetPointsValue :one
SELECT points
FROM point_sources
WHERE source = @event
AND guild_id = @guild_id;

