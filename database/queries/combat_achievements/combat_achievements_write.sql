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

-- name: GiveUserCombatAchievement :exec
INSERT INTO user_combat_achievement (user_id, guild_id, combat_achievement_name)
VALUES (@user_id, @guild_id, @combat_achievement_name)
ON CONFLICT ON CONSTRAINT "user_combat_achievement_pkey" DO NOTHING;

-- name: RemoveUserCombatAchievement :execrows
DELETE FROM user_combat_achievement
WHERE user_id = @user_id
AND guild_id = @guild_id
AND combat_achievement_name = @combat_achievement_name;
