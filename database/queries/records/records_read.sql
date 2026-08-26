-- name: GetBosses :many
SELECT name, display_name, category, solo, value_type FROM bosses;

-- name: GetCategories :many
SELECT "thumbnail", "order", "name" FROM categories;

-- name: GetUserRecords :many
SELECT
    r.record_id,
    r.boss_name,
    b.display_name,
    b.category,
    b.solo,
    b.value_type,
    r.date,
    r.value,
    tm.user_id,
    tm.guild_id
FROM records r
JOIN teams tm ON r.record_id = tm.record_id AND r.guild_id = tm.guild_id
JOIN bosses b ON r.boss_name = b.name
WHERE tm.user_id = @user_id AND tm.guild_id = @guild_id
ORDER BY r.record_id;

-- name: GetBossInfo :one
SELECT b.name, b.display_name, b.category, b.solo, b.value_type, vt.higher_is_better
FROM bosses b
JOIN value_types vt ON b.value_type = vt.name
WHERE b.name = @boss_name;

-- name: GetBossRecords :many
SELECT r.record_id, r.value, r.boss_name, r.date, r.guild_id, tm.user_id
FROM records r
JOIN teams tm ON r.record_id = tm.record_id AND r.guild_id = tm.guild_id
WHERE r.guild_id = @guild_id AND r.boss_name = @boss_name
ORDER BY r.record_id, tm.user_id;

