-- name: UpdatePointsByEvent :many
UPDATE users
SET points = (
	SELECT points
    FROM point_sources
    WHERE source = $1
    AND "point_sources"."guild_id" = $2
)
WHERE user_id = ANY($3::text[])
AND guild_id = $2 RETURNING *;

-- name: UpdatePointsCustom :many
UPDATE users
SET points = points + 10
WHERE user_id = ANY($1::text[])
AND guild_id = $1 RETURNING *;
