-- +goose Up
-- +goose StatementBegin

CREATE TABLE "public"."guild_ranks" (
    "guild_id" character varying(32) NOT NULL,
    "name" character varying(32) NOT NULL,
    "min_points" integer NOT NULL DEFAULT 0,
    "icon" text,
    "role_id" character varying(32),
    "display_order" smallint NOT NULL DEFAULT 0,
    CONSTRAINT "guild_ranks_pkey" PRIMARY KEY ("guild_id", "name")
) WITH (oids = false);

ALTER TABLE "public"."guild_ranks"
ADD CONSTRAINT "guild_ranks_guild_id_fkey" FOREIGN KEY ("guild_id")
REFERENCES "guilds" ("guild_id") ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

-- Trigger function to seed default guild ranks on guild creation
CREATE OR REPLACE FUNCTION insert_default_guild_ranks()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO guild_ranks (guild_id, name, min_points, display_order)
  VALUES
    (NEW.guild_id, 'jade',        0,    1),
    (NEW.guild_id, 'red_topaz',   50,   2),
    (NEW.guild_id, 'sapphire',    100,  3),
    (NEW.guild_id, 'emerald',     200,  4),
    (NEW.guild_id, 'ruby',        400,  5),
    (NEW.guild_id, 'diamond',     600,  6),
    (NEW.guild_id, 'dragonstone', 800,  7),
    (NEW.guild_id, 'onyx',        1000, 8),
    (NEW.guild_id, 'zenyte',      1250, 9),
    (NEW.guild_id, 'astral',      1500, 10),
    (NEW.guild_id, 'death',       2000, 11),
    (NEW.guild_id, 'blood',       2750, 12),
    (NEW.guild_id, 'soul',        3750, 13),
    (NEW.guild_id, 'wrath',       5000, 14)
  ON CONFLICT ON CONSTRAINT guild_ranks_pkey DO NOTHING;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER seed_default_guild_ranks_trigger
AFTER INSERT ON guilds
FOR EACH ROW
EXECUTE FUNCTION insert_default_guild_ranks();

-- Seed default ranks for all existing guilds
INSERT INTO guild_ranks (guild_id, name, min_points, display_order)
SELECT g.guild_id, r.name, r.min_points, r.display_order
FROM guilds g
CROSS JOIN (VALUES
    ('jade',        0,    1),
    ('red_topaz',   50,   2),
    ('sapphire',    100,  3),
    ('emerald',     200,  4),
    ('ruby',        400,  5),
    ('diamond',     600,  6),
    ('dragonstone', 800,  7),
    ('onyx',        1000, 8),
    ('zenyte',      1250, 9),
    ('astral',      1500, 10),
    ('death',       2000, 11),
    ('blood',       2750, 12),
    ('soul',        3750, 13),
    ('wrath',       5000, 14)
) AS r(name, min_points, display_order)
ON CONFLICT ON CONSTRAINT guild_ranks_pkey DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS seed_default_guild_ranks_trigger ON guilds;
DROP FUNCTION IF EXISTS insert_default_guild_ranks();
DROP TABLE IF EXISTS "guild_ranks";
-- +goose StatementEnd
