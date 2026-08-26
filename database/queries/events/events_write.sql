-- name: CreateEvent :exec
INSERT INTO event (
	name,
	wom_id,
	guild_id,
	position_cutoff,
	solo
) VALUES (
	@name,
	@wom_id,
	@guild_id,
	@position_cutoff,
	@solo
);

-- name: InsertEventParticipants :exec
WITH participant_data AS (
    SELECT
        unnest(@participant_ids::text[]) as wom_id,
		generate_series(1, ARRAY_LENGTH(@participant_ids::text[], 1)) as placement
)
INSERT INTO event_participant (
    user_id,
    placement,
    guild_id,
    event_id
)
SELECT
    r.user_id,
    pd.placement,
    @guild_id,
    @wom_id
FROM participant_data pd
JOIN rsn r ON r.wom_id = pd.wom_id AND r.guild_id = @guild_id;

-- name: InsertEventTeams :exec
WITH participant_data AS (
    SELECT
        unnest(@participant_ids::text[]) as wom_id,
        unnest(@participant_placements::int[]) as placement
)
INSERT INTO event_participant (
    user_id,
    guild_id,
    placement,
    event_id
)
SELECT
    r.user_id,
    @guild_id,
    pd.placement,
    @wom_id
FROM participant_data pd
JOIN rsn r ON r.wom_id = pd.wom_id AND r.guild_id = @guild_id;

-- name: InsertLegacyEventParticipants :exec
INSERT INTO event_participant (
    user_id,
    placement,
    guild_id,
    event_id
)
SELECT
    unnest(@user_ids::text[]),
    unnest(@placements::int[]),
    @guild_id,
    @event_id
ON CONFLICT DO NOTHING;

-- name: DeleteEvent :exec
DELETE FROM event WHERE wom_id = @event_id;

-- name: UpdatePointsByEvent :many
WITH point_value AS (
    SELECT points
    FROM point_sources
    WHERE source = @event
    AND guild_id = @guild_id
)
UPDATE users
SET points = points + (SELECT points FROM point_value)
WHERE user_id = ANY(@user_ids::text[])
AND users.guild_id = @guild_id
RETURNING user_id, guild_id, points, (SELECT points FROM point_value) AS given_points;

-- name: UpdateEvent :one
UPDATE event SET
    name = COALESCE(sqlc.narg('name'), name),
    position_cutoff = COALESCE(sqlc.narg('position_cutoff'), position_cutoff),
    solo = COALESCE(sqlc.narg('solo'), solo)
WHERE guild_id = @guild_id AND wom_id = @wom_id
RETURNING name, guild_id, wom_id, position_cutoff, solo;
