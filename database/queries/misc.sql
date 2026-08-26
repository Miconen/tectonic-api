-- name: GetLeaderboard :many
SELECT u.user_id, u.guild_id, u.points, json_agg(r) AS rsns
FROM users u
JOIN rsn r ON u.user_id = r.user_id AND u.guild_id = r.guild_id
WHERE u.guild_id = @guild_id
GROUP BY u.user_id, u.guild_id, u.points
ORDER BY u.points DESC
LIMIT @user_limit;

