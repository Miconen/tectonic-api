-- +goose Up
-- +goose StatementBegin
CREATE VIEW detailed_event AS
SELECT
    e.name,
    e.wom_id,
    e.guild_id,
    e.position_cutoff,
    ep.user_id,
    ep.placement
FROM event_participant ep, event e;

DROP VIEW detailed_users;

CREATE VIEW detailed_users AS
SELECT
    u.user_id,
    u.guild_id,
    u.points,
    array_remove(array_agg(r), NULL) AS rsns,
    array_remove(array_agg(de), NULL) as events
FROM users u
JOIN rsn r ON u.user_id = r.user_id AND u.guild_id = r.guild_id
LEFT JOIN detailed_event de ON u.user_id = de.user_id AND u.guild_id = de.guild_id
GROUP BY u.user_id, u.guild_id, u.points;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW detailed_users;
DROP VIEW detailed_event;

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
