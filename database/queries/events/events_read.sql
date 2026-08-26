-- name: GetGuildEvents :many
SELECT "name", "wom_id", "guild_id", "position_cutoff", "solo" FROM event WHERE guild_id = @guild_id;

-- name: GetEventParticipation :many
SELECT
	ep.user_id,
	ep.placement
FROM event_participant ep WHERE ep.event_id = @event_id;

-- name: GetUserEvents :many
SELECT
    e.name,
    e.wom_id AS event_id,
    e.guild_id,
    ep.user_id,
    ep.placement,
    e.position_cutoff,
    e.solo
FROM event e
JOIN event_participant ep ON e.wom_id = ep.event_id
WHERE ep.user_id = @user_id AND ep.guild_id = @guild_id
AND ep.placement <= e.position_cutoff;

