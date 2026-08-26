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

-- name: GetGuildRanks :many
SELECT name, min_points, icon, role_id, display_order
FROM guild_ranks
WHERE guild_id = @guild_id
ORDER BY display_order;
