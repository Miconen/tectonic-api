-- name: CreateGuild :one
INSERT INTO guilds (
  guild_id
) VALUES (
  $1
)
RETURNING guild_id, multiplier, pb_channel_id, position_count;

-- name: DeleteGuild :execrows
DELETE FROM guilds
WHERE guild_id = $1;

-- name: GetGuild :one
SELECT
    guilds.guild_id, guilds.multiplier, guilds.pb_channel_id, guilds.mod_channel_id, guilds.log_channel_id, guilds.position_count,
    (SELECT count(user_id) FROM users WHERE users.guild_id = $1) as user_count,
    (SELECT count(record_id) FROM records WHERE records.guild_id = $1) as record_count
FROM guilds
WHERE guilds.guild_id = $1 LIMIT 1;

-- name: UpdateGuild :one
UPDATE guilds SET
    multiplier = COALESCE(sqlc.narg('multiplier'), multiplier),
    pb_channel_id = COALESCE(sqlc.narg('pb_channel_id'), pb_channel_id),
    mod_channel_id = COALESCE(sqlc.narg('mod_channel_id'), mod_channel_id),
    log_channel_id = COALESCE(sqlc.narg('log_channel_id'), log_channel_id),
    position_count = COALESCE(sqlc.narg('position_count'), position_count)
WHERE guild_id = @guild_id
RETURNING guild_id, multiplier, pb_channel_id;
