-- name: CreateGuildRank :exec
INSERT INTO guild_ranks (guild_id, name, min_points, icon, role_id, display_order)
VALUES (@guild_id, @name, @min_points, @icon, @role_id, @display_order);

-- name: UpdateGuildRank :execrows
UPDATE guild_ranks SET
    min_points = COALESCE(sqlc.narg('min_points'), min_points),
    icon = COALESCE(sqlc.narg('icon'), icon),
    role_id = COALESCE(sqlc.narg('role_id'), role_id),
    display_order = COALESCE(sqlc.narg('display_order'), display_order)
WHERE guild_id = @guild_id AND name = @name;


-- name: DeleteGuildRank :execrows
DELETE FROM guild_ranks WHERE guild_id = @guild_id AND name = @name;
