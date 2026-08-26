-- name: GetDetailedUserBases :many
WITH ranked_users AS (
    SELECT
        user_id,
        guild_id,
        RANK() OVER (
            PARTITION BY guild_id
            ORDER BY points DESC
        ) AS user_rank
    FROM users
    WHERE guild_id = @guild_id
)
SELECT
    u.user_id,
    u.guild_id,
    u.points,
    ru.user_rank,
    tier.name AS tier_name,
    tier.icon AS tier_icon,
    tier.role_id AS tier_role_id,
    tier.min_points AS tier_min_points,
    tier.display_order AS tier_display_order
FROM users u
JOIN ranked_users ru
  ON ru.user_id = u.user_id
 AND ru.guild_id = u.guild_id
LEFT JOIN LATERAL (
    SELECT
        gr.name,
        gr.icon,
        gr.role_id,
        gr.min_points,
        gr.display_order
    FROM guild_ranks gr
    WHERE gr.guild_id = u.guild_id
      AND gr.min_points <= u.points
    ORDER BY gr.min_points DESC
    LIMIT 1
) tier ON TRUE
WHERE u.guild_id = @guild_id
  AND u.user_id = ANY(@user_ids::text[]);

-- name: GetDetailedUserRsns :many
SELECT
    r.user_id,
    r.guild_id,
    r.rsn,
    r.wom_id
FROM rsn r
WHERE r.guild_id = @guild_id
  AND r.user_id = ANY(@user_ids::text[])
ORDER BY r.user_id, r.rsn;

-- name: GetDetailedUserAchievements :many
SELECT
    ua.user_id,
    ua.guild_id,
    a.name,
    a.thumbnail,
    a.discord_icon,
    a.order AS achievement_order
FROM user_achievement ua
JOIN achievement a
  ON a.name = ua.achievement_name
WHERE ua.guild_id = @guild_id
  AND ua.user_id = ANY(@user_ids::text[])
ORDER BY ua.user_id, a.order;

-- name: GetDetailedUserCombatAchievements :many
SELECT
    uca.user_id,
    uca.guild_id,
    uca.combat_achievement_name AS name
FROM user_combat_achievement uca
WHERE uca.guild_id = @guild_id
  AND uca.user_id = ANY(@user_ids::text[])
ORDER BY uca.user_id, uca.combat_achievement_name;

-- name: GetDetailedUserEvents :many
SELECT
    ep.user_id,
    ep.guild_id,
    e.name,
    e.wom_id AS event_id,
    ep.placement,
    e.position_cutoff,
    e.solo
FROM event_participant ep
JOIN event e
  ON e.wom_id = ep.event_id
 AND e.guild_id = ep.guild_id
WHERE ep.guild_id = @guild_id
  AND ep.user_id = ANY(@user_ids::text[])
  AND ep.placement <= e.position_cutoff
ORDER BY ep.user_id, e.wom_id;

-- name: GetDetailedUserRecords :many
WITH selected_records AS (
    SELECT DISTINCT
        tm.user_id AS owner_user_id,
        tm.guild_id,
        tm.record_id
    FROM teams tm
    WHERE tm.guild_id = @guild_id
      AND tm.user_id = ANY(@user_ids::text[])
)
SELECT
    sr.owner_user_id,
    sr.guild_id,
    r.record_id,
    r.boss_name,
    b.display_name,
    b.category,
    b.solo,
    b.value_type,
    r.date,
    r.value,
    teammate.user_id AS teammate_user_id,
    teammate.guild_id AS teammate_guild_id
FROM selected_records sr
JOIN records r
  ON r.record_id = sr.record_id
 AND r.guild_id = sr.guild_id
JOIN bosses b
  ON b.name = r.boss_name
JOIN teams teammate
  ON teammate.record_id = r.record_id
 AND teammate.guild_id = r.guild_id
ORDER BY
    sr.owner_user_id,
    r.record_id,
    teammate.user_id;
