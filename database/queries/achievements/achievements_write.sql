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
