-- +goose Up
-- +goose StatementBegin
CREATE VIEW detailed_time_teams AS
SELECT
    tm.run_id,
    tm.user_id,
    tm.guild_id,
    u.points,
    array_remove(array_agg(r), NULL) AS rsns
FROM teams tm
JOIN users u ON tm.user_id = u.user_id
JOIN rsn r ON tm.user_id = r.user_id AND u.guild_id = r.guild_id
GROUP BY tm.user_id, tm.guild_id, u.points, tm.run_id;

CREATE VIEW detailed_times AS
SELECT
    t.time,
    t.boss_name,
    b.display_name,
    b.category,
    t.run_id,
    t.date,
    array_remove(array_agg(dtt), NULL) AS team
FROM times t
LEFT JOIN detailed_time_teams dtt ON t.run_id = dtt.run_id
LEFT JOIN bosses b ON b.name = t.boss_name
GROUP BY t.time, t.boss_name, b.category, b.display_name, t.run_id, t.date;

CREATE VIEW detailed_users AS
SELECT
    u.user_id,
    u.guild_id,
    u.points,
    array_remove(array_agg(r), NULL) AS rsns
FROM users u
JOIN rsn r ON u.user_id = r.user_id AND u.guild_id = r.guild_id
GROUP BY u.user_id, u.guild_id, u.points;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW detailed_users;
DROP VIEW detailed_times;
DROP VIEW detailed_time_teams;
-- +goose StatementEnd
