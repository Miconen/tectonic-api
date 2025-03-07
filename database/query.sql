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
WITH query_user AS (
    SELECT u.user_id, u.guild_id, u.points,
    array_agg(r) AS rsns
    FROM users u, rsn r
    WHERE u.user_id = r.user_id
    GROUP BY u.user_id, u.guild_id, u.points
), user_times AS (
    SELECT
        t.time,
        t.boss_name,
        b.category,
        t.run_id,
        t.date
    FROM times t, bosses b
    WHERE b.name = t.boss_name
), time_teammates AS (
    SELECT 
        ut.run_id,
        tm.user_id,
        tm.guild_id,
        u.points,
        array_agg(r) AS rsns
    FROM teams tm, users u, rsn r, user_times ut
    WHERE tm.run_id = ut.run_id
    AND tm.user_id = r.user_id
    GROUP BY tm.user_id, tm.guild_id, u.points, ut.run_id
), time_with_teammates AS (
    SELECT
        ut.time,
        ut.boss_name,
        ut.category,
        ut.run_id,
        ut.date,
        array_remove(array_agg(tt), NULL) AS teammates
    FROM user_times ut
    LEFT JOIN time_teammates tt ON ut.run_id = tt.run_id
    GROUP BY ut.time, ut.boss_name, ut.category, ut.run_id, ut.date
)
SELECT
    qu.user_id,
    qu.guild_id,
    qu.points,
    to_json(qu.rsns) AS rsns,
    json_agg(twt) AS times
FROM query_user qu, time_with_teammates twt
GROUP BY qu.user_id, qu.guild_id, qu.points, qu.rsns;

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
