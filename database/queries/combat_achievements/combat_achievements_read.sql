-- name: GetCombatAchievement :one
SELECT ca.name, ca.point_source
FROM combat_achievement ca
WHERE ca.name = @name AND ca.guild_id = @guild_id;

-- name: GetUserCombatAchievements :many
SELECT uca.combat_achievement_name
FROM user_combat_achievement uca
WHERE uca.user_id = @user_id AND uca.guild_id = @guild_id;

-- name: GetGuildCombatAchievements :many
SELECT ca.name, ca.point_source, ps.points, ps.name AS point_source_display_name
FROM combat_achievement ca
JOIN point_sources ps ON ca.guild_id = ps.guild_id AND ca.point_source = ps.source
WHERE ca.guild_id = @guild_id
ORDER BY ca.name;

