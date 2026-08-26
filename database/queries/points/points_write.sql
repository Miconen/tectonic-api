-- name: UpdatePointsCustom :many
UPDATE users
SET points = points + @points
WHERE user_id = ANY(@user_ids::text[])
AND guild_id = @guild_id
RETURNING user_id, guild_id, points, @points::int AS given_points;

-- name: UpdateGuildPointSource :execrows
UPDATE point_sources ps
SET points = @points
WHERE ps.guild_id = @guild_id
AND ps.source = @point_source;
