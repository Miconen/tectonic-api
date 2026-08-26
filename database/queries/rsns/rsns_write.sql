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

-- name: DeleteRsn :execrows
DELETE FROM rsn r
WHERE r.guild_id = @guild_id AND r.user_id = @user_id AND r.rsn = @rsn;
