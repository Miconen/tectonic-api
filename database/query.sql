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
RETURNING guild_id, multiplier, pb_channel_id, position_count;

-- name: DeleteGuild :execrows
DELETE FROM guilds
WHERE guild_id = $1;

-- name: GetGuild :one
SELECT
    guilds.guild_id, guilds.multiplier, guilds.pb_channel_id, guilds.mod_channel_id, guilds.position_count,
    (SELECT count(user_id) FROM users WHERE users.guild_id = $1) as user_count,
    (SELECT count(record_id) FROM records WHERE records.guild_id = $1) as record_count
FROM guilds
WHERE guilds.guild_id = $1 LIMIT 1;

-- name: UpdateGuild :one
UPDATE guilds SET
    multiplier = CASE WHEN @multiplier::numeric IS NOT NULL AND @multiplier::numeric != 0 THEN @multiplier::numeric ELSE multiplier END,
    pb_channel_id = CASE WHEN @pb_channel_id::text IS NOT NULL AND @pb_channel_id::text != '' THEN @pb_channel_id::text ELSE pb_channel_id END,
    mod_channel_id = CASE WHEN @mod_channel_id::text IS NOT NULL AND @mod_channel_id::text != '' THEN @mod_channel_id::text ELSE mod_channel_id END,
    position_count = CASE WHEN @position_count::smallint IS NOT NULL AND @position_count::smallint != 0 THEN @position_count::smallint ELSE position_count END
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

-- ==================== Teams ====================

-- name: AddToTeamByBoss :exec
INSERT INTO teams (record_id, user_id, guild_id)
VALUES (
    (SELECT r.record_id
     FROM records r
     JOIN bosses b ON r.boss_name = b.name
     JOIN value_types vt ON b.value_type = vt.name
     WHERE r.guild_id = @guild_id
       AND r.boss_name = @boss_name
     ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
     LIMIT 1),
    @user_id,
    @guild_id
);

-- name: AddToTeamById :exec
INSERT INTO teams (record_id, user_id, guild_id)
VALUES (
    @record_id,
    @user_id,
    @guild_id
);

-- name: RemoveFromTeamByBoss :execrows
DELETE FROM teams
WHERE record_id = (
    SELECT r.record_id
    FROM records r
    JOIN bosses b ON r.boss_name = b.name
    JOIN value_types vt ON b.value_type = vt.name
    WHERE r.guild_id = @guild_id
      AND r.boss_name = @boss_name
    ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
    LIMIT 1
)
AND user_id = @user_id
AND guild_id = @guild_id;

-- name: RemoveFromTeamById :execrows
DELETE FROM teams
WHERE record_id = @record_id
AND user_id = @user_id
AND guild_id = @guild_id;

-- name: CreateTeam :exec
INSERT INTO teams (record_id, user_id, guild_id)
SELECT @record_id, unnest(@user_ids::text[]), @guild_id;

-- ==================== Records ====================

-- name: CreateRecord :one
INSERT INTO records (
    value,
    boss_name,
    date,
    guild_id
)
VALUES (
    @value,
    @boss_name,
    @date,
    @guild_id
) RETURNING record_id;

-- name: DeleteRecord :execrows
DELETE FROM records r
WHERE r.guild_id = @guild_id AND r.record_id = @record_id;

-- name: DeleteBossRecords :execrows
DELETE FROM records
WHERE guild_id = @guild_id AND boss_name = @boss_name;

-- name: DeleteTopRecord :execrows
DELETE FROM records
WHERE record_id = (
    SELECT r.record_id
    FROM records r
    JOIN bosses b ON r.boss_name = b.name
    JOIN value_types vt ON b.value_type = vt.name
    WHERE r.guild_id = @guild_id AND r.boss_name = @boss_name
    ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
    LIMIT 1
)
AND guild_id = @guild_id;

-- name: GetBossInfo :one
SELECT b.name, b.display_name, b.category, b.solo, b.value_type, vt.higher_is_better
FROM bosses b
JOIN value_types vt ON b.value_type = vt.name
WHERE b.name = @boss_name;

-- name: GetBossRecords :many
SELECT r.record_id, r.value, r.boss_name, r.date, r.guild_id, tm.user_id
FROM records r
JOIN teams tm ON r.record_id = tm.record_id AND r.guild_id = tm.guild_id
WHERE r.guild_id = @guild_id AND r.boss_name = @boss_name
ORDER BY r.record_id, tm.user_id;

-- ==================== Detailed Guild ====================

-- name: GetDetailedGuild :one
WITH ranked_records AS (
    -- Team boss records ranked directly
    SELECT r.record_id, r.value, r.boss_name, r.date, r.guild_id,
           ROW_NUMBER() OVER (
               PARTITION BY r.boss_name
               ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
           ) as position
    FROM records r
    JOIN bosses b ON r.boss_name = b.name
    JOIN value_types vt ON b.value_type = vt.name
    WHERE r.guild_id = @guild_id AND b.solo = false

    UNION ALL

    -- Solo boss records: best per user, then ranked
    SELECT s.record_id, s.value, s.boss_name, s.date, s.guild_id,
           ROW_NUMBER() OVER (
               PARTITION BY s.boss_name
               ORDER BY CASE WHEN s.higher_is_better THEN -s.value ELSE s.value END ASC
           ) as position
    FROM (
        SELECT DISTINCT ON (tm.user_id, r.boss_name)
            r.record_id, r.value, r.boss_name, r.date, r.guild_id, vt.higher_is_better
        FROM records r
        JOIN teams tm ON r.record_id = tm.record_id AND r.guild_id = tm.guild_id
        JOIN bosses b ON r.boss_name = b.name
        JOIN value_types vt ON b.value_type = vt.name
        WHERE r.guild_id = @guild_id AND b.solo = true
        ORDER BY tm.user_id, r.boss_name,
                 CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
    ) s
),
top_records AS (
    SELECT record_id, value, boss_name, date, guild_id, position
    FROM ranked_records
    WHERE position <= (SELECT position_count FROM guilds WHERE guild_id = @guild_id)
)
SELECT
    g.guild_id,
    g.multiplier,
    g.pb_channel_id,
    g.mod_channel_id,
    g.position_count,
    (SELECT count(user_id) FROM users WHERE users.guild_id = @guild_id) as user_count,
    (SELECT count(record_id) FROM records WHERE records.guild_id = @guild_id) as record_count,

    (SELECT json_agg(tm) FROM teams tm
     WHERE tm.guild_id = g.guild_id
     AND tm.record_id IN (SELECT tr.record_id FROM top_records tr)) AS teammates,

    (SELECT json_agg(tr) FROM top_records tr) AS records,

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

-- ==================== User Records ====================

-- name: GetUserRecords :many
SELECT
    r.record_id,
    r.boss_name,
    b.display_name,
    b.category,
    b.solo,
    b.value_type,
    r.date,
    r.value,
    tm.user_id,
    tm.guild_id
FROM records r
JOIN teams tm ON r.record_id = tm.record_id AND r.guild_id = tm.guild_id
JOIN bosses b ON r.boss_name = b.name
WHERE tm.user_id = @user_id AND tm.guild_id = @guild_id
ORDER BY r.record_id;

-- ==================== User Rank & Tier ====================

-- name: GetUserRank :one
WITH ranked_users AS (
    SELECT user_id, RANK() OVER (ORDER BY points DESC) as user_rank
    FROM users
    WHERE guild_id = @guild_id
)
SELECT user_rank FROM ranked_users
WHERE user_id = @user_id;

-- name: GetUserTier :one
SELECT gr.name, gr.icon, gr.role_id, gr.min_points, gr.display_order
FROM guild_ranks gr
WHERE gr.guild_id = @guild_id
AND gr.min_points <= @points
ORDER BY gr.min_points DESC
LIMIT 1;

-- ==================== Guild Ranks ====================

-- name: GetGuildRanks :many
SELECT name, min_points, icon, role_id, display_order
FROM guild_ranks
WHERE guild_id = @guild_id
ORDER BY display_order;

-- name: CreateGuildRank :exec
INSERT INTO guild_ranks (guild_id, name, min_points, icon, role_id, display_order)
VALUES (@guild_id, @name, @min_points, @icon, @role_id, @display_order);

-- name: UpdateGuildRank :execrows
UPDATE guild_ranks SET
    min_points = CASE WHEN @min_points::int IS NOT NULL AND @min_points::int != -1 THEN @min_points::int ELSE min_points END,
    icon = CASE WHEN @icon::text IS NOT NULL AND @icon::text != '' THEN @icon::text ELSE icon END,
    role_id = CASE WHEN @role_id::text IS NOT NULL AND @role_id::text != '' THEN @role_id::text ELSE role_id END,
    display_order = CASE WHEN @display_order::smallint IS NOT NULL AND @display_order::smallint != -1 THEN @display_order::smallint ELSE display_order END
WHERE guild_id = @guild_id AND name = @name;

-- name: DeleteGuildRank :execrows
DELETE FROM guild_ranks WHERE guild_id = @guild_id AND name = @name;

-- ==================== Value Types ====================

-- name: GetValueTypes :many
SELECT name, higher_is_better FROM value_types ORDER BY name;

-- ==================== Misc ====================

-- name: UpdateCategoryMessageIds :execrows
UPDATE guild_categories
SET message_id = u.message_id
FROM (SELECT unnest(@categories::text[]) as category,
             unnest(@message_ids::text[]) as message_id) as u
WHERE guild_categories.guild_id = @guild_id
AND guild_categories.category = u.category;

-- name: GetBosses :many
SELECT name, display_name, category, solo, value_type FROM bosses;

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

-- name: GetGuildCombatAchievements :many
SELECT ca.name, ca.point_source, ps.points, ps.name AS point_source_display_name
FROM combat_achievement ca
JOIN point_sources ps ON ca.guild_id = ps.guild_id AND ca.point_source = ps.source
WHERE ca.guild_id = @guild_id
ORDER BY ca.name;

-- name: GetCombatAchievement :one
SELECT ca.name, ca.point_source
FROM combat_achievement ca
WHERE ca.name = @name AND ca.guild_id = @guild_id;

-- name: CreateCombatAchievement :exec
INSERT INTO combat_achievement (name, guild_id, point_source)
VALUES (@name, @guild_id, @point_source);

-- name: DeleteCombatAchievement :execrows
DELETE FROM combat_achievement
WHERE name = @name AND guild_id = @guild_id;

-- name: CompleteCombatAchievement :exec
INSERT INTO user_combat_achievement (user_id, guild_id, combat_achievement_name)
SELECT unnest(@user_ids::text[]), @guild_id, @combat_achievement_name
ON CONFLICT ON CONSTRAINT "user_combat_achievement_pkey" DO NOTHING;

-- name: GetUserCombatAchievements :many
SELECT uca.combat_achievement_name
FROM user_combat_achievement uca
WHERE uca.user_id = @user_id AND uca.guild_id = @guild_id;

-- name: GiveUserCombatAchievement :exec
INSERT INTO user_combat_achievement (user_id, guild_id, combat_achievement_name)
VALUES (@user_id, @guild_id, @combat_achievement_name)
ON CONFLICT ON CONSTRAINT "user_combat_achievement_pkey" DO NOTHING;

-- name: RemoveUserCombatAchievement :execrows
DELETE FROM user_combat_achievement
WHERE user_id = @user_id
AND guild_id = @guild_id
AND combat_achievement_name = @combat_achievement_name;

-- name: DeleteRsn :execrows
DELETE FROM rsn r
WHERE r.guild_id = @guild_id AND r.user_id = @user_id AND r.rsn = @rsn;
