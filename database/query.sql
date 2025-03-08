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
SET points = points + @points
WHERE user_id = ANY(@user_ids::text[])
AND guild_id = @guild_id RETURNING user_id, guild_id, points;

-- name: GetDetailedUsers :many
SELECT 
    du.user_id,
    du.guild_id,
    du.points,
    to_json(du.rsns) AS rsns,
    COALESCE(times_json, '[]'::json) AS times
FROM detailed_users du
LEFT JOIN LATERAL (
    SELECT json_agg(dt) AS times_json
    FROM detailed_times dt
    WHERE dt.run_id IN (
        SELECT tm.run_id
        FROM teams tm
        WHERE tm.user_id = du.user_id AND tm.guild_id = du.guild_id
    )
) t ON true
WHERE du.user_id = ANY(@user_ids::text[])
AND du.guild_id = @guild_id;

-- name: GetDetailedUsersByRSN :many
WITH rsn_user AS (
    SELECT r.user_id, r.guild_id 
    FROM rsn r 
    WHERE r.rsn = ANY(@rsns::text[])
)
SELECT 
    du.user_id,
    du.guild_id,
    du.points,
    to_json(du.rsns) AS rsns,
    COALESCE(times_json, '[]'::json) AS times
FROM detailed_users du
JOIN rsn_user ru ON ru.user_id = du.user_id AND ru.guild_id = du.guild_id
LEFT JOIN LATERAL (
    SELECT json_agg(dt) AS times_json
    FROM detailed_times dt
    WHERE dt.run_id IN (
        SELECT tm.run_id
        FROM teams tm
        WHERE tm.user_id = du.user_id AND tm.guild_id = du.guild_id
    )
) t ON true
WHERE du.guild_id = @guild_id;

-- name: GetDetailedUsersByWomID :many
WITH wom_user AS (
    SELECT r.user_id, r.guild_id 
    FROM rsn r 
    WHERE r.wom_id = ANY(@wom_ids::text[])
)
SELECT 
    du.user_id,
    du.guild_id,
    du.points,
    to_json(du.rsns) AS rsns,
    COALESCE(times_json, '[]'::json) AS times
FROM detailed_users du
JOIN wom_user wu ON wu.user_id = du.user_id AND wu.guild_id = du.guild_id
LEFT JOIN LATERAL (
    SELECT json_agg(dt) AS times_json
    FROM detailed_times dt
    WHERE dt.run_id IN (
        SELECT tm.run_id
        FROM teams tm
        WHERE tm.user_id = du.user_id AND tm.guild_id = du.guild_id
    )
) t ON true
WHERE du.guild_id = @guild_id;

-- name: GetLeaderboard :many
SELECT u.user_id, u.guild_id, u.points, json_agg(r) AS rsns
FROM users u
JOIN rsn r ON u.user_id = r.user_id AND u.guild_id = r.guild_id
WHERE u.guild_id = @guild_id
GROUP BY u.user_id, u.guild_id, u.points
ORDER BY u.points DESC
LIMIT @user_limit;

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
