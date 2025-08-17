-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."combat_achievement" (
    "name" character varying(32) NOT NULL,
    "guild_id" character varying(32) NOT NULL,
    "point_source" character varying(32) NOT NULL,
    CONSTRAINT "combat_achievement_pkey" PRIMARY KEY ("name", "guild_id")
) WITH (oids = false);

ALTER TABLE "public"."combat_achievement" ADD FOREIGN KEY (guild_id, point_source) REFERENCES "point_sources" (guild_id, source) ON DELETE CASCADE ON UPDATE CASCADE;

INSERT INTO "point_sources" (guild_id, source, points)
SELECT g.guild_id, 'combat_achievement_low', 5
FROM guilds g;

INSERT INTO "point_sources" (guild_id, source, points)
SELECT g.guild_id, 'combat_achievement_medium', 10
FROM guilds g;

INSERT INTO "point_sources" (guild_id, source, points)
SELECT g.guild_id, 'combat_achievement_high', 15
FROM guilds g;

INSERT INTO "point_sources" (guild_id, source, points)
SELECT g.guild_id, 'combat_achievement_very_high', 25
FROM guilds g;

INSERT INTO "point_sources" (guild_id, source, points)
SELECT g.guild_id, 'combat_achievement_highest', 50
FROM guilds g;

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

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "combat_achievement";

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
    (NEW.guild_id, 'split_high', 30)
  ON CONFLICT ON CONSTRAINT point_sources_pkey DO NOTHING;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DELETE FROM "point_sources"
WHERE source = 'combat_achievement_low'
OR 'combat_achievement_medium'
OR 'combat_achievement_high'
OR 'combat_achievement_very_high'
OR 'combat_achievement_highest';
-- +goose StatementEnd
