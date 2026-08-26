-- name: AddToTeamByBoss :exec
INSERT INTO teams (record_id, user_id, guild_id)
VALUES (
    (SELECT r.record_id
     FROM records r
     JOIN bosses b ON r.boss_name = b.name
     JOIN value_types vt ON b.value_type = vt.name
     WHERE r.guild_id = @guild_id
       AND r.boss_name = @boss_name
     ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
     LIMIT 1),
    @user_id,
    @guild_id
);

-- name: AddToTeamById :exec
INSERT INTO teams (record_id, user_id, guild_id)
VALUES (
    @record_id,
    @user_id,
    @guild_id
);

-- name: RemoveFromTeamByBoss :execrows
DELETE FROM teams
WHERE record_id = (
    SELECT r.record_id
    FROM records r
    JOIN bosses b ON r.boss_name = b.name
    JOIN value_types vt ON b.value_type = vt.name
    WHERE r.guild_id = @guild_id
      AND r.boss_name = @boss_name
    ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
    LIMIT 1
)
AND user_id = @user_id
AND guild_id = @guild_id;

-- name: RemoveFromTeamById :execrows
DELETE FROM teams
WHERE record_id = @record_id
AND user_id = @user_id
AND guild_id = @guild_id;

-- name: CreateTeam :exec
INSERT INTO teams (record_id, user_id, guild_id)
SELECT @record_id, unnest(@user_ids::text[]), @guild_id;

