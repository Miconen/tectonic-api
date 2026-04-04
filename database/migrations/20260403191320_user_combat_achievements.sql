-- +goose Up
-- +goose StatementBegin

-- Widen name column for longer CA names
ALTER TABLE "combat_achievement"
ALTER COLUMN "name" TYPE character varying(64);

-- Track which users have completed which CAs
CREATE TABLE "public"."user_combat_achievement" (
    "user_id" character varying(32) NOT NULL,
    "guild_id" character varying(32) NOT NULL,
    "combat_achievement_name" character varying(64) NOT NULL,
    CONSTRAINT "user_combat_achievement_pkey" PRIMARY KEY ("user_id", "guild_id", "combat_achievement_name")
) WITH (oids = false);

ALTER TABLE "public"."user_combat_achievement"
    ADD CONSTRAINT "user_combat_achievement_user_fkey"
    FOREIGN KEY ("user_id", "guild_id") REFERENCES "users"("user_id", "guild_id")
    ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "public"."user_combat_achievement"
    ADD CONSTRAINT "user_combat_achievement_ca_fkey"
    FOREIGN KEY ("combat_achievement_name", "guild_id") REFERENCES "combat_achievement"("name", "guild_id")
    ON DELETE CASCADE ON UPDATE CASCADE;

-- Seed default combat achievements for all existing guilds
INSERT INTO "combat_achievement" (name, guild_id, point_source)
SELECT ca.name, g.guild_id, ca.point_source
FROM guilds g
CROSS JOIN (VALUES
    -- Tombs of Amascut Kits
    ('Ava''s kit', 'combat_achievement_low'),
    ('Ward kit', 'combat_achievement_low'),
    ('Fang kit', 'combat_achievement_medium'),
    ('Kephri transmog', 'combat_achievement_low'),
    ('Akkha transmog', 'combat_achievement_medium'),
    ('Ba-ba transmog', 'combat_achievement_low'),
    ('Zebak transmog', 'combat_achievement_low'),
    ('Warden transmog', 'combat_achievement_low'),
    -- Tombs of Amascut CAs
    ('Perfect Warden', 'combat_achievement_medium'),
    ('Perfection of Scabaras', 'combat_achievement_low'),
    ('Perfection of Crondis', 'combat_achievement_low'),
    ('Perfection of Apmeken', 'combat_achievement_low'),
    ('Perfection of Het', 'combat_achievement_low'),
    ('Insanity', 'combat_achievement_high'),
    ('All Praise Zebak', 'combat_achievement_medium'),
    ('Fancy Feet', 'combat_achievement_low'),
    ('But... Damage', 'combat_achievement_low'),
    ('All Out of Medics', 'combat_achievement_low'),
    ('Tombs Speed Runner', 'combat_achievement_medium'),
    ('Tombs Speed Runner II', 'combat_achievement_high'),
    ('Tombs Speed Runner III', 'combat_achievement_very_high'),
    -- Chambers of Xeric CAs
    ('Trio CoX Speedrunner', 'combat_achievement_low'),
    ('Five-Man CoX Speedrunner', 'combat_achievement_low'),
    ('Trio CoX Perfect Olm', 'combat_achievement_medium'),
    -- Chambers of Xeric: CM CAs
    ('Trio CoX:CM Speedrunner', 'combat_achievement_medium'),
    ('Five-Man CoX:CM Speedrunner', 'combat_achievement_medium'),
    -- Theatre of Blood CAs
    ('Duo ToB Speedrunner', 'combat_achievement_medium'),
    ('Trio ToB Speedrunner', 'combat_achievement_medium'),
    ('Four-Man ToB Speedrunner', 'combat_achievement_medium'),
    ('Five-Man ToB Speedrunner', 'combat_achievement_medium'),
    ('Morytania only', 'combat_achievement_medium'),
    ('Back in my day', 'combat_achievement_low'),
    ('Can''t drain this', 'combat_achievement_low'),
    ('Pop it', 'combat_achievement_medium'),
    ('Perfect Theatre', 'combat_achievement_highest'),
    ('Perfect Verzik', 'combat_achievement_high'),
    ('Perfect Bloat', 'combat_achievement_low'),
    ('Perfect Maiden', 'combat_achievement_low'),
    ('Perfect Nylocas', 'combat_achievement_low'),
    ('Perfect Soteseg', 'combat_achievement_low'),
    ('Perfect Xarpus', 'combat_achievement_low'),
    -- Theatre of Blood: HM CAs
    ('Trio HMT Speedrunner', 'combat_achievement_low'),
    ('Four-Man HMT Speedrunner', 'combat_achievement_low'),
    ('Five-Man HMT Speedrunner', 'combat_achievement_low'),
    ('Pack Like a Yak', 'combat_achievement_low'),
    ('Team work makes the dream work', 'combat_achievement_low'),
    ('Royal Affairs', 'combat_achievement_medium'),
    ('Nylo Sniper', 'combat_achievement_low'),
    ('Personal Space', 'combat_achievement_medium'),
    ('Harder Mode I', 'combat_achievement_low'),
    ('Harder Mode II', 'combat_achievement_high'),
    ('Harder Mode III', 'combat_achievement_high'),
    -- Nex
    ('Nex Duo', 'combat_achievement_medium'),
    -- Nightmare
    ('Five-Man Nightmare Speed-Chaser', 'combat_achievement_low'),
    ('Five-Man Nightmare Speed-Runner', 'combat_achievement_medium'),
    ('A Long Trip', 'combat_achievement_medium'),
    -- Yama
    ('Contract Choreographer', 'combat_achievement_low'),
    ('Yama Speed-Runner', 'combat_achievement_low'),
    ('Contractually Unbound', 'combat_achievement_low'),
    -- Huey
    ('Hueycoatl Speed-Runner', 'combat_achievement_low')
) AS ca(name, point_source)
ON CONFLICT ON CONSTRAINT "combat_achievement_pkey" DO NOTHING;

-- Separate trigger function for seeding combat achievements on new guilds
CREATE OR REPLACE FUNCTION insert_default_combat_achievements()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO combat_achievement (name, guild_id, point_source)
  VALUES
    -- Tombs of Amascut Kits
    ('Ava''s kit', NEW.guild_id, 'combat_achievement_low'),
    ('Ward kit', NEW.guild_id, 'combat_achievement_low'),
    ('Fang kit', NEW.guild_id, 'combat_achievement_medium'),
    ('Kephri transmog', NEW.guild_id, 'combat_achievement_low'),
    ('Akkha transmog', NEW.guild_id, 'combat_achievement_medium'),
    ('Ba-ba transmog', NEW.guild_id, 'combat_achievement_low'),
    ('Zebak transmog', NEW.guild_id, 'combat_achievement_low'),
    ('Warden transmog', NEW.guild_id, 'combat_achievement_low'),
    -- Tombs of Amascut CAs
    ('Perfect Warden', NEW.guild_id, 'combat_achievement_medium'),
    ('Perfection of Scabaras', NEW.guild_id, 'combat_achievement_low'),
    ('Perfection of Crondis', NEW.guild_id, 'combat_achievement_low'),
    ('Perfection of Apmeken', NEW.guild_id, 'combat_achievement_low'),
    ('Perfection of Het', NEW.guild_id, 'combat_achievement_low'),
    ('Insanity', NEW.guild_id, 'combat_achievement_high'),
    ('All Praise Zebak', NEW.guild_id, 'combat_achievement_medium'),
    ('Fancy Feet', NEW.guild_id, 'combat_achievement_low'),
    ('But... Damage', NEW.guild_id, 'combat_achievement_low'),
    ('All Out of Medics', NEW.guild_id, 'combat_achievement_low'),
    ('Tombs Speed Runner', NEW.guild_id, 'combat_achievement_medium'),
    ('Tombs Speed Runner II', NEW.guild_id, 'combat_achievement_high'),
    ('Tombs Speed Runner III', NEW.guild_id, 'combat_achievement_very_high'),
    -- Chambers of Xeric CAs
    ('Trio CoX Speedrunner', NEW.guild_id, 'combat_achievement_low'),
    ('Five-Man CoX Speedrunner', NEW.guild_id, 'combat_achievement_low'),
    ('Trio CoX Perfect Olm', NEW.guild_id, 'combat_achievement_medium'),
    -- Chambers of Xeric: CM CAs
    ('Trio CoX:CM Speedrunner', NEW.guild_id, 'combat_achievement_medium'),
    ('Five-Man CoX:CM Speedrunner', NEW.guild_id, 'combat_achievement_medium'),
    -- Theatre of Blood CAs
    ('Duo ToB Speedrunner', NEW.guild_id, 'combat_achievement_medium'),
    ('Trio ToB Speedrunner', NEW.guild_id, 'combat_achievement_medium'),
    ('Four-Man ToB Speedrunner', NEW.guild_id, 'combat_achievement_medium'),
    ('Five-Man ToB Speedrunner', NEW.guild_id, 'combat_achievement_medium'),
    ('Morytania only', NEW.guild_id, 'combat_achievement_medium'),
    ('Back in my day', NEW.guild_id, 'combat_achievement_low'),
    ('Can''t drain this', NEW.guild_id, 'combat_achievement_low'),
    ('Pop it', NEW.guild_id, 'combat_achievement_medium'),
    ('Perfect Theatre', NEW.guild_id, 'combat_achievement_highest'),
    ('Perfect Verzik', NEW.guild_id, 'combat_achievement_high'),
    ('Perfect Bloat', NEW.guild_id, 'combat_achievement_low'),
    ('Perfect Maiden', NEW.guild_id, 'combat_achievement_low'),
    ('Perfect Nylocas', NEW.guild_id, 'combat_achievement_low'),
    ('Perfect Soteseg', NEW.guild_id, 'combat_achievement_low'),
    ('Perfect Xarpus', NEW.guild_id, 'combat_achievement_low'),
    -- Theatre of Blood: HM CAs
    ('Trio HMT Speedrunner', NEW.guild_id, 'combat_achievement_low'),
    ('Four-Man HMT Speedrunner', NEW.guild_id, 'combat_achievement_low'),
    ('Five-Man HMT Speedrunner', NEW.guild_id, 'combat_achievement_low'),
    ('Pack Like a Yak', NEW.guild_id, 'combat_achievement_low'),
    ('Team work makes the dream work', NEW.guild_id, 'combat_achievement_low'),
    ('Royal Affairs', NEW.guild_id, 'combat_achievement_medium'),
    ('Nylo Sniper', NEW.guild_id, 'combat_achievement_low'),
    ('Personal Space', NEW.guild_id, 'combat_achievement_medium'),
    ('Harder Mode I', NEW.guild_id, 'combat_achievement_low'),
    ('Harder Mode II', NEW.guild_id, 'combat_achievement_high'),
    ('Harder Mode III', NEW.guild_id, 'combat_achievement_high'),
    -- Nex
    ('Nex Duo', NEW.guild_id, 'combat_achievement_medium'),
    -- Nightmare
    ('Five-Man Nightmare Speed-Chaser', NEW.guild_id, 'combat_achievement_low'),
    ('Five-Man Nightmare Speed-Runner', NEW.guild_id, 'combat_achievement_medium'),
    ('A Long Trip', NEW.guild_id, 'combat_achievement_medium'),
    -- Yama
    ('Contract Choreographer', NEW.guild_id, 'combat_achievement_low'),
    ('Yama Speed-Runner', NEW.guild_id, 'combat_achievement_low'),
    ('Contractually Unbound', NEW.guild_id, 'combat_achievement_low'),
    -- Huey
    ('Hueycoatl Speed-Runner', NEW.guild_id, 'combat_achievement_low')
  ON CONFLICT ON CONSTRAINT "combat_achievement_pkey" DO NOTHING;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_default_combat_achievements_trigger
AFTER INSERT ON guilds
FOR EACH ROW
EXECUTE FUNCTION insert_default_combat_achievements();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS insert_default_combat_achievements_trigger ON guilds;
DROP FUNCTION IF EXISTS insert_default_combat_achievements();

DELETE FROM "combat_achievement"
WHERE name IN (
    'Ava''s kit', 'Ward kit', 'Fang kit', 'Kephri transmog', 'Akkha transmog',
    'Ba-ba transmog', 'Zebak transmog', 'Warden transmog', 'Perfect Warden',
    'Perfection of Scabaras', 'Perfection of Crondis', 'Perfection of Apmeken',
    'Perfection of Het', 'Insanity', 'All Praise Zebak', 'Fancy Feet',
    'But... Damage', 'All Out of Medics', 'Tombs Speed Runner',
    'Tombs Speed Runner II', 'Tombs Speed Runner III', 'Trio CoX Speedrunner',
    'Five-Man CoX Speedrunner', 'Trio CoX Perfect Olm', 'Trio CoX:CM Speedrunner',
    'Five-Man CoX:CM Speedrunner', 'Duo ToB Speedrunner', 'Trio ToB Speedrunner',
    'Four-Man ToB Speedrunner', 'Five-Man ToB Speedrunner', 'Morytania only',
    'Back in my day', 'Can''t drain this', 'Pop it', 'Perfect Theatre',
    'Perfect Verzik', 'Perfect Bloat', 'Perfect Maiden', 'Perfect Nylocas',
    'Perfect Soteseg', 'Perfect Xarpus', 'Trio HMT Speedrunner',
    'Four-Man HMT Speedrunner', 'Five-Man HMT Speedrunner', 'Pack Like a Yak',
    'Team work makes the dream work', 'Royal Affairs', 'Nylo Sniper',
    'Personal Space', 'Harder Mode I', 'Harder Mode II', 'Harder Mode III',
    'Nex Duo', 'Five-Man Nightmare Speed-Chaser', 'Five-Man Nightmare Speed-Runner',
    'A Long Trip', 'Contract Choreographer', 'Yama Speed-Runner',
    'Contractually Unbound', 'Hueycoatl Speed-Runner'
);

DROP TABLE IF EXISTS "user_combat_achievement";

ALTER TABLE "combat_achievement"
ALTER COLUMN "name" TYPE character varying(32);

-- +goose StatementEnd
