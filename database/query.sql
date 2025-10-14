-- name: GetUsersById :many
select users.user_id, users.guild_id, users.points
from users
where users.guild_id = $1
and users.user_id = any(@user_ids::text[]);

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
AND guild_id = @guild_id
RETURNING user_id, guild_id, points, @points::int AS given_points;

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

-- name: DeleteGuild :execrows
DELETE FROM guilds
WHERE guild_id = $1;

-- name: GetGuild :one
SELECT guild_id, multiplier, pb_channel_id FROM guilds
WHERE guild_id = $1 LIMIT 1;

-- name: UpdateGuild :one
UPDATE guilds SET
    multiplier = CASE WHEN @multiplier::numeric IS NOT NULL AND @multiplier::numeric != 0 THEN @multiplier::numeric ELSE multiplier END,
    pb_channel_id = CASE WHEN @pb_channel_id::text IS NOT NULL AND @pb_channel_id::text != '' THEN @pb_channel_id::text ELSE pb_channel_id END
WHERE guild_id = @guild_id RETURNING guild_id, multiplier, pb_channel_id;

-- name: UpdateEvent :one
UPDATE event SET
    name = CASE WHEN @name::text IS NOT NULL AND @name::text != '' THEN @name::text ELSE name END,
    position_cutoff = CASE WHEN @position_cutoff::numeric IS NOT NULL AND @position_cutoff::numeric != 0 THEN @position_cutoff::numeric ELSE position_cutoff END,
    solo = CASE WHEN @solo::boolean IS NOT NULL THEN @solo::boolean ELSE solo END
WHERE guild_id = @guild_id
AND wom_id = @wom_id
RETURNING name, guild_id, wom_id, position_cutoff;

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

-- name: AddToTeamByBoss :exec
INSERT INTO teams (run_id, user_id, guild_id) 
VALUES (
    (SELECT pb_id 
     FROM guild_bosses 
     WHERE guild_id = @guild_id 
       AND boss = @boss_name),
    @user_id,
    @guild_id
);

-- name: AddToTeamById :exec
INSERT INTO teams (run_id, user_id, guild_id) 
VALUES (
    @run_id,
    @user_id,
    @guild_id
);

-- name: RemoveFromTeamByBoss :execrows
DELETE FROM teams 
WHERE run_id = (
    SELECT pb_id 
    FROM guild_bosses 
     WHERE guild_bosses.guild_id = @guild_id
       AND boss = @boss_name
)
AND user_id = @user_id
AND guild_id = @guild_id;

-- name: RemoveFromTeamById :execrows
DELETE FROM teams 
WHERE run_id = @run_id
AND user_id = @user_id
AND guild_id = @guild_id;

-- name: DeleteRsn :execrows
DELETE FROM rsn r
WHERE r.guild_id = @guild_id AND r.user_id = @user_id AND r.rsn = @rsn;

-- name: DeleteTime :execrows
DELETE FROM times t
WHERE t.guild_id = @guild_id AND t.run_id = @run_id;

-- name: DeletePb :execrows
DELETE FROM times t
WHERE t.guild_id = @guild_id
AND t.boss_name = @boss_name
AND t.run_id = (
    SELECT pb_id
    FROM guild_bosses
    WHERE guild_id = @guild_id
    AND boss = @boss_name
);

-- name: RevertGuildBossPb :exec
UPDATE guild_bosses 
SET pb_id = (
    SELECT run_id 
    FROM times 
    WHERE times.guild_id = @guild_id
      AND times.boss_name = @boss_name
    ORDER BY time ASC 
    LIMIT 1
)
WHERE guild_bosses.guild_id = @guild_id
AND guild_bosses.boss = @boss_name;

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
     JOIN guild_bosses gb ON t.run_id = gb.pb_id
     WHERE t.guild_id = g.guild_id
     	AND gb.guild_id = g.guild_id) AS teammates,

    (SELECT json_agg(t) FROM times t 
     JOIN guild_bosses gb ON t.run_id = gb.pb_id 
     WHERE t.guild_id = g.guild_id
     	AND gb.guild_id = g.guild_id
     	AND gb.pb_id IS NOT NULL) AS pbs,

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

-- name: GetAchievements :many
SELECT "name", "thumbnail", "discord_icon", "order" FROM achievement;

-- name: GetGuildEvents :many
SELECT "name", "wom_id", "guild_id", "position_cutoff", "solo" FROM event WHERE guild_id = @guild_id;

-- name: GetGuildPointSources :many
SELECT "source", "points", "name" FROM point_sources WHERE guild_id = @guild_id;

-- name: UpdateGuildPointSource :execrows
UPDATE point_sources ps
SET points = @points
WHERE ps.guild_id = @guild_id
AND ps.source = @point_source;

-- name: CreateEvent :exec
INSERT INTO event (
	name,
	wom_id,
	guild_id,
	position_cutoff,
	solo
) VALUES (
	@name,
	@wom_id,
	@guild_id,
	@position_cutoff,
	@solo
);

-- name: InsertEventParticipants :exec
WITH participant_data AS (
    SELECT 
        unnest(@participant_ids::text[]) as wom_id,
		generate_series(1, ARRAY_LENGTH(@participant_ids::text[], 1)) as placement
)
INSERT INTO event_participant (
    user_id,
    placement,
    guild_id,
    event_id
) 
SELECT 
    r.user_id,
    pd.placement,
    @guild_id,
    @wom_id
FROM participant_data pd
JOIN rsn r ON r.wom_id = pd.wom_id AND r.guild_id = @guild_id;

-- name: InsertEventTeams :exec
WITH participant_data AS (
    SELECT 
        unnest(@participant_ids::text[]) as wom_id,
        unnest(@participant_placements::int[]) as placement
)
INSERT INTO event_participant (
    user_id,
    guild_id,
    placement,
    event_id
) 
SELECT 
    r.user_id,
    @guild_id,
    pd.placement,
    @wom_id
FROM participant_data pd
JOIN rsn r ON r.wom_id = pd.wom_id AND r.guild_id = @guild_id;

-- name: DeleteEvent :exec
DELETE FROM event WHERE wom_id = @event_id;

-- name: GetEventParticipation :many
SELECT
	ep.user_id,
	ep.placement
FROM event_participant ep WHERE ep.event_id = @event_id;

-- name: GetUserAchievements :many
SELECT
	a.name,
	a.thumbnail,
	a.discord_icon
FROM user_achievement ua
JOIN achievement a ON ua.achievement_name = a.name
WHERE ua.user_id = @user_id
ORDER BY a.order;

-- name: GetUserTimes :many
SELECT
    t.run_id,
    t.boss_name,
    b.display_name,
    b.category,
    b.solo,
    t.date,
    t.time,
    tm.user_id,
    tm.guild_id
FROM times t
JOIN teams tm ON t.run_id = tm.run_id
JOIN guild_bosses gb ON t.run_id = gb.pb_id AND tm.guild_id = gb.guild_id
JOIN bosses b ON b.name = gb.boss
WHERE tm.user_id = @user_id AND tm.guild_id = @guild_id AND gb.pb_id = t.run_id
ORDER BY t.run_id;

-- name: GetUserRsns :many
SELECT
	r.rsn,
	r.wom_id
FROM rsn r
WHERE r.user_id = @user_id AND r.guild_id = @guild_id;

-- name: GetUserByWom :many
SELECT
	r.user_id
FROM rsn r
WHERE r.wom_id = ANY(@wom_id::text[]);

-- name: GetGuildUserByWom :many
SELECT
	r.user_id
FROM rsn r
WHERE r.wom_id = ANY(@wom_ids::text[])
AND r.guild_id = @guild_id;

-- name: GetGuildUserByRsn :many
SELECT
	r.user_id
FROM rsn r
WHERE r.rsn ILIKE ANY(@rsns::text[])
AND r.guild_id = @guild_id;

-- name: GetUserByRsn :many
SELECT
	r.user_id
FROM rsn r
WHERE r.rsn = ANY(@rsns::text[]);

-- name: GetUserEvents :many
SELECT
    e.name,
    e.wom_id AS event_id,
    e.guild_id,
    ep.user_id,
    ep.placement,
    e.position_cutoff,
    e.solo
FROM event e
JOIN event_participant ep ON e.wom_id = ep.event_id
WHERE ep.user_id = @user_id AND ep.guild_id = @guild_id
AND ep.placement <= e.position_cutoff;

-- name: GiveAchievementById :exec
INSERT INTO user_achievement (
	user_id,
	achievement_name,
	guild_id
) VALUES (
	@user_id,
	@achievement_name,
	@guild_id
);

-- name: GiveAchievementByRsn :exec
WITH user_lookup AS (
    SELECT user_id FROM rsn WHERE rsn = @rsn AND guild_id = @guild_id
)
INSERT INTO user_achievement (
	user_id,
	achievement_name,
	guild_id
) SELECT user_id, @achievement_name, @guild_id
FROM user_lookup;

-- name: RemoveAchievementById :exec
DELETE FROM user_achievement ua
WHERE ua.user_id = @user_id
AND ua.achievement_name = @achievement_name
AND ua.guild_id = @guild_id;

-- name: RemoveAchievementByRsn :exec
DELETE FROM user_achievement ua
WHERE ua.user_id IN (SELECT r.user_id FROM rsn r WHERE r.rsn = @rsn)
AND ua.achievement_name = @achievement_name
AND ua.guild_id = @guild_id;
