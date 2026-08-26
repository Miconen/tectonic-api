-- name: GetUserRsns :many
SELECT
	r.rsn,
	r.wom_id
FROM rsn r
WHERE r.user_id = @user_id AND r.guild_id = @guild_id;

