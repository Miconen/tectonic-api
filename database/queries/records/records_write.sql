-- name: GetValueTypes :many
SELECT name, higher_is_better FROM value_types ORDER BY name;

-- name: CreateRecord :one
INSERT INTO records (
    value,
    boss_name,
    date,
    guild_id
)
VALUES (
    @value,
    @boss_name,
    @date,
    @guild_id
) RETURNING record_id;

-- name: DeleteRecord :execrows
DELETE FROM records r
WHERE r.guild_id = @guild_id AND r.record_id = @record_id;

-- name: DeleteRecordsByUserId :execrows
DELETE FROM records r
WHERE r.guild_id = @guild_id
AND record_id IN (
    SELECT t.record_id
    FROM teams t
    WHERE t.guild_id = @guild_id
      AND t.user_id = @user_id
);

-- name: DeleteRecordsByRsn :execrows
DELETE FROM records r
WHERE r.guild_id = @guild_id
AND record_id IN (
    SELECT t.record_id
    FROM teams t
    WHERE t.guild_id = @guild_id
      AND t.user_id IN (
          SELECT r.user_id
          FROM rsn r
          WHERE r.guild_id = @guild_id AND r.rsn = @rsn
      )
);

-- name: DeleteRecordsByWom :execrows
DELETE FROM records r
WHERE r.guild_id = @guild_id
AND record_id IN (
    SELECT t.record_id
    FROM teams t
    WHERE t.guild_id = @guild_id
      AND t.user_id IN (
          SELECT r.user_id
          FROM rsn r
          WHERE r.guild_id = @guild_id AND r.wom_id = @wom_id
      )
);

-- name: DeleteBossRecords :execrows
DELETE FROM records
WHERE guild_id = @guild_id AND boss_name = @boss_name;

-- name: DeleteTopRecord :execrows
DELETE FROM records
WHERE record_id = (
    SELECT r.record_id
    FROM records r
    JOIN bosses b ON r.boss_name = b.name
    JOIN value_types vt ON b.value_type = vt.name
    WHERE r.guild_id = @guild_id AND r.boss_name = @boss_name
    ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
    LIMIT 1
)
AND guild_id = @guild_id;

-- name: UpdateCategoryMessageIds :execrows
UPDATE guild_categories
SET message_id = u.message_id
FROM (SELECT unnest(@categories::text[]) as category,
             unnest(@message_ids::text[]) as message_id) as u
WHERE guild_categories.guild_id = @guild_id
AND guild_categories.category = u.category;

