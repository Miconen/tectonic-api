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

-- name: DeleteUserById :execrows
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id = $2;

-- name: DeleteUserByRsn :execrows
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.rsn = $2
);

-- name: DeleteUserByWom :execrows
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.wom_id = $2
);

-- name: GetPointsValue :one
SELECT points
FROM point_sources
WHERE source = @event
AND guild_id = @guild_id;

-- name: UpdatePointsByEvent :many
WITH point_value AS (
    SELECT points
    FROM point_sources
    WHERE source = @event
    AND guild_id = @guild_id
)
UPDATE users
SET points = points + (SELECT points FROM point_value)
WHERE user_id = ANY(@user_ids::text[])
AND users.guild_id = @guild_id 
RETURNING user_id, guild_id, points, (SELECT points FROM point_value) AS given_points;

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
    multiplier = CASE WHEN @multiplier::numeric IS NOT NULL AND @multiplier::numeric != 0 THEN @multiplier::numeric ELSE multiplier END,
    pb_channel_id = CASE WHEN @pb_channel_id::text IS NOT NULL AND @pb_channel_id::text != '' THEN @pb_channel_id::text ELSE pb_channel_id END
WHERE guild_id = @guild_id RETURNING guild_id, multiplier, pb_channel_id;

-- name: CreateRsn :exec
INSERT INTO rsn (
    guild_id,
    user_id,
    rsn,
    wom_id
) VALUES (
    @guild_id,
    @user_id,
    @rsn,
    @wom_id
);

-- name: DeleteRsn :execrows
DELETE FROM rsn r
WHERE r.guild_id = @guild_id AND r.user_id = @user_id AND r.rsn = @rsn;

-- name: DeleteTime :execrows
DELETE FROM times t
WHERE t.guild_id = @guild_id AND t.run_id = @run_id;

-- name: CreateTime :one
INSERT INTO times (
    time,
    boss_name,
    date,
    guild_id
)
VALUES (
    @time,
    @boss_name,
    @date,
    @guild_id
) RETURNING run_id;

-- name: CheckPb :one
SELECT gb.boss, t.time
FROM guild_bosses gb
LEFT JOIN times t ON gb.pb_id = t.run_id AND gb.guild_id = t.guild_id
WHERE gb.boss = @boss
AND gb.guild_id = @guild_id;

-- name: UpdatePb :execrows
UPDATE guild_bosses SET
    pb_id = @run_id
WHERE guild_id = @guild_id
AND boss = @boss;

-- name: CreateTeam :exec
INSERT INTO teams (run_id, user_id, guild_id)
SELECT @run_id, unnest(@user_ids::text[]), @guild_id;

-- name: GetDetailedGuild :one
SELECT
    g.guild_id,
    g.pb_channel_id,
    (SELECT json_agg(tm) FROM teams tm 
     JOIN times t ON tm.run_id = t.run_id 
     WHERE t.guild_id = g.guild_id) AS teammates,
    (SELECT json_agg(t) FROM times t 
     WHERE t.guild_id = g.guild_id) AS pbs,
    (SELECT json_agg(b) FROM bosses b 
     JOIN guild_bosses gb ON b.name = gb.boss 
     WHERE gb.guild_id = g.guild_id) AS bosses,
    (SELECT json_agg(c) FROM categories c 
     JOIN guild_categories gc ON c.name = gc.category 
     WHERE gc.guild_id = g.guild_id) AS categories,
    (SELECT json_agg(gb) FROM guild_bosses gb 
     WHERE gb.guild_id = g.guild_id) AS guild_bosses,
    (SELECT json_agg(gc) FROM guild_categories gc 
     WHERE gc.guild_id = g.guild_id) AS guild_categories
FROM guilds g
WHERE g.guild_id = @guild_id; 

-- name: UpdateCategoryMessageIds :execrows
UPDATE guild_categories
SET message_id = u.message_id
FROM (SELECT unnest(@categories::text[]) as category,
             unnest(@message_ids::text[]) as message_id) as u
WHERE guild_categories.guild_id = @guild_id
AND guild_categories.category = u.category;

-- name: GetBosses :many
SELECT name, display_name, category, solo FROM bosses;

-- name: GetCategories :many
SELECT "thumbnail", "order", "name" FROM categories;
