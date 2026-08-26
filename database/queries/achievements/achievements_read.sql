-- name: GetAchievements :many
SELECT "name", "thumbnail", "discord_icon", "order" FROM achievement;

-- name: GetUserAchievements :many
SELECT
	a.name,
	a.thumbnail,
	a.discord_icon
FROM user_achievement ua
JOIN achievement a ON ua.achievement_name = a.name
WHERE ua.user_id = @user_id
ORDER BY a.order;

