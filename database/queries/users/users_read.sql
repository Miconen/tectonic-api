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
