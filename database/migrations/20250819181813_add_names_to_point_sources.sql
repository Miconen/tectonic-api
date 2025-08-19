-- +goose Up
-- +goose StatementBegin
ALTER TABLE "point_sources"
ADD COLUMN "name" character varying(64);

UPDATE "point_sources"
SET "name" = point_sources.source
WHERE "name" IS NULL;

ALTER TABLE "point_sources"
ALTER COLUMN "name" SET NOT NULL;

CREATE OR REPLACE FUNCTION insert_default_point_sources()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO point_sources (guild_id, source, points, name)
  VALUES
    (NEW.guild_id, 'event_participation', 5, 'Event participation'),
    (NEW.guild_id, 'event_hosting', 10, 'Event hosting'),
    (NEW.guild_id, 'clan_pb', 10, 'Clan PB'),
    (NEW.guild_id, 'split_low', 10, 'Split (Low)'),
    (NEW.guild_id, 'split_medium', 20, 'Split (Medium)'),
    (NEW.guild_id, 'split_high', 30, 'Split (High)'),
    (NEW.guild_id, 'combat_achievement_low', 5, 'Combat Achievement (Low)'),
    (NEW.guild_id, 'combat_achievement_medium', 10, 'Combat Achievement (Medium)'),
    (NEW.guild_id, 'combat_achievement_high', 15, 'Combat Achievement (High)'),
    (NEW.guild_id, 'combat_achievement_very_high', 25, 'Combat Achievement (Very High)'),
    (NEW.guild_id, 'combat_achievement_highest', 50, 'Combat Achievement (Highest)')
  ON CONFLICT ON CONSTRAINT point_sources_pkey DO NOTHING;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "point_sources"
DROP COLUMN "name";

CREATE OR REPLACE FUNCTION insert_default_point_sources()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO point_sources (guild_id, source, points)
  VALUES
    (NEW.guild_id, 'event_participation', 5),
    (NEW.guild_id, 'event_hosting', 10),
    (NEW.guild_id, 'clan_pb', 10),
    (NEW.guild_id, 'split_low', 10),
    (NEW.guild_id, 'split_medium', 20),
    (NEW.guild_id, 'split_high', 30),
    (NEW.guild_id, 'combat_achievement_low', 5),
    (NEW.guild_id, 'combat_achievement_medium', 10),
    (NEW.guild_id, 'combat_achievement_high', 15),
    (NEW.guild_id, 'combat_achievement_very_high', 25),
    (NEW.guild_id, 'combat_achievement_highest', 50)
  ON CONFLICT ON CONSTRAINT point_sources_pkey DO NOTHING;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
