-- name: GetUsersById :many
SELECT users.user_id, users.guild_id, users.points
FROM users
WHERE users.guild_id = $1
AND users.user_id = ANY(@user_ids::text[]);

-- name: GetUsersByRsn :many
SELECT users.user_id, users.guild_id, users.points
FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.rsn = ANY(@rsns::text[])
);

-- name: GetUsersByWom :many
SELECT users.user_id, users.guild_id, users.points
FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.wom_id = ANY(@wom_ids::text[])
);

-- name: CreateUser :one
WITH inserted_users AS (
  INSERT INTO users (guild_id, user_id)
  VALUES ($1, $2)
  RETURNING guild_id, user_id, points
),
inserted_rsn AS (
  INSERT INTO rsn (guild_id, user_id, rsn, wom_id)
  VALUES ($1, $2, $3, $4)
  RETURNING guild_id, user_id, rsn, wom_id
)
SELECT
    u.guild_id,
    u.user_id,
    u.points,
    r.rsn,
    r.wom_id
FROM inserted_users u
JOIN inserted_rsn r
ON u.guild_id = r.guild_id AND u.user_id = r.user_id;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id = $2;

-- name: DeleteUserByRsn :exec
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.rsn = $2
);

-- name: DeleteUserByWom :exec
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.wom_id = $2
);

-- name: UpdatePointsByEvent :many
UPDATE users
SET points = points + (
	SELECT points
    FROM point_sources
    WHERE source = @event
    AND "point_sources"."guild_id" = @guild_id
)
WHERE user_id = ANY(@user_ids::text[])
AND guild_id = @guild_id RETURNING user_id, guild_id, points;

-- name: UpdatePointsCustom :many
UPDATE users
SET points = points + $1
WHERE user_id = ANY($2::text[])
AND guild_id = $3 RETURNING user_id, guild_id, points;

-- name: GetDetailedUsers :many
WITH guild_pb_runs AS (
    SELECT gb.guild_id, gb.pb_id
    FROM guild_bosses gb
)
SELECT 
    u.user_id AS user_id,
    u.guild_id AS user_guild_id,
    u.points AS user_points,
    r.rsn AS user_rsn,
    r.wom_id AS user_wom_id,
    t.time AS time_value,
    t.boss_name,
    b.category AS boss_category,
    t.run_id,
    t.date AS run_date,
    tm.user_id AS team_user_id,
    tm.guild_id AS team_guild_id,
    tu.points AS team_user_points,
    tr.rsn AS team_user_rsn,
    tr.wom_id AS team_user_wom_id
FROM users u
LEFT JOIN rsn r ON r.user_id = u.user_id AND r.guild_id = u.guild_id
LEFT JOIN guild_pb_runs g ON g.guild_id = u.guild_id
LEFT JOIN times t ON t.run_id = g.pb_id
LEFT JOIN bosses b ON b.name = t.boss_name
LEFT JOIN teams tm ON tm.run_id = t.run_id
LEFT JOIN users tu ON tm.user_id = tu.user_id AND tm.guild_id = tu.guild_id
LEFT JOIN rsn tr ON tr.user_id = tu.user_id AND tr.guild_id = tu.guild_id
ORDER BY
    u.user_id,
    u.guild_id,
    t.run_id,
    tm.user_id;

-- name: CreateGuild :one
INSERT INTO guilds (
  guild_id
) VALUES (
  $1
)
RETURNING guild_id, multiplier, pb_channel_id;


-- name: DeleteGuild :one
DELETE FROM guilds
WHERE guild_id = $1 RETURNING guild_id, multiplier, pb_channel_id;

-- name: GetGuild :one
SELECT guild_id, multiplier, pb_channel_id FROM guilds
WHERE guild_id = $1 LIMIT 1;

-- name: UpdateGuild :one
UPDATE guilds SET
    multiplier = CASE WHEN $1::numeric IS NOT NULL AND $1::numeric != 0 THEN $1::numeric ELSE multiplier END,
    pb_channel_id = CASE WHEN $2::text IS NOT NULL AND $2::text != '' THEN $2::text ELSE pb_channel_id END
WHERE guild_id = $3 RETURNING guild_id, multiplier, pb_channel_id;
