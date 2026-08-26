-- name: GetDetailedGuild :one
WITH ranked_records AS (
    -- Team boss records ranked directly
    SELECT r.record_id, r.value, r.boss_name, r.date, r.guild_id,
           ROW_NUMBER() OVER (
               PARTITION BY r.boss_name
               ORDER BY CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
           ) as position
    FROM records r
    JOIN bosses b ON r.boss_name = b.name
    JOIN value_types vt ON b.value_type = vt.name
    WHERE r.guild_id = @guild_id AND b.solo = false

    UNION ALL

    -- Solo boss records: best per user, then ranked
    SELECT s.record_id, s.value, s.boss_name, s.date, s.guild_id,
           ROW_NUMBER() OVER (
               PARTITION BY s.boss_name
               ORDER BY CASE WHEN s.higher_is_better THEN -s.value ELSE s.value END ASC
           ) as position
    FROM (
        SELECT DISTINCT ON (tm.user_id, r.boss_name)
            r.record_id, r.value, r.boss_name, r.date, r.guild_id, vt.higher_is_better
        FROM records r
        JOIN teams tm ON r.record_id = tm.record_id AND r.guild_id = tm.guild_id
        JOIN bosses b ON r.boss_name = b.name
        JOIN value_types vt ON b.value_type = vt.name
        WHERE r.guild_id = @guild_id AND b.solo = true
        ORDER BY tm.user_id, r.boss_name,
                 CASE WHEN vt.higher_is_better THEN -r.value ELSE r.value END ASC
    ) s
),
top_records AS (
    SELECT record_id, value, boss_name, date, guild_id, position
    FROM ranked_records
    WHERE position <= (SELECT position_count FROM guilds WHERE guild_id = @guild_id)
)
SELECT
    g.guild_id,
    g.multiplier,
    g.pb_channel_id,
    g.mod_channel_id,
    g.log_channel_id,
    g.position_count,
    (SELECT count(user_id) FROM users WHERE users.guild_id = @guild_id) as user_count,
    (SELECT count(record_id) FROM records WHERE records.guild_id = @guild_id) as record_count,

    (SELECT json_agg(tm) FROM teams tm
     WHERE tm.guild_id = g.guild_id
     AND tm.record_id IN (SELECT tr.record_id FROM top_records tr)) AS teammates,

    (SELECT json_agg(tr) FROM top_records tr) AS records,

    (SELECT json_agg(b) FROM bosses b
     JOIN guild_bosses gb ON b.name = gb.boss
     WHERE gb.guild_id = g.guild_id) AS bosses,

    (SELECT json_agg(c) FROM categories c
     JOIN guild_categories gc ON c.name = gc.category
     WHERE gc.guild_id = g.guild_id) AS categories,

    (SELECT json_agg(gb) FROM guild_bosses gb
     WHERE gb.guild_id = g.guild_id) AS guild_bosses,

    (SELECT json_agg(gc) FROM guild_categories gc
     WHERE gc.guild_id = g.guild_id) AS guild_categories

FROM guilds g
WHERE g.guild_id = @guild_id;
