-- +goose Up
-- +goose StatementBegin

-- Add missing boss PBs (Maggot King, Mad Angel)
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('maggot_king', 'Maggot King', 'Miscellaneous', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('mad_angel', 'Mad Angel', 'Miscellaneous', '1');

INSERT INTO "guild_bosses" ("boss", "guild_id", "category")
SELECT boss_list.boss, g.guild_id, 'Miscellaneous'
FROM "guilds" g
CROSS JOIN (VALUES ('maggot_king'), ('mad_angel')) AS boss_list(boss);

-- Remove the any-scale CoX/CoX:CM PBs and the ToA 400 invocation tier
-- (cascades to guild_bosses and any existing records for these bosses)
DELETE FROM "bosses"
WHERE "name" IN ('cox_any', 'cm_any', 'toa_solo_400', 'toa_team_400');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('cox_any', 'Any', 'Chambers of Xeric', '0');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('cm_any', 'Any', 'Chambers of Xeric: CM', '0');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('toa_solo_400', 'Solo 400+', 'Tombs of Amascut', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('toa_team_400', 'Team 400+', 'Tombs of Amascut', '0');

INSERT INTO "guild_bosses" ("boss", "guild_id", "category")
SELECT boss_list.boss, g.guild_id, boss_list.category
FROM "guilds" g
CROSS JOIN (VALUES
    ('cox_any', 'Chambers of Xeric'),
    ('cm_any', 'Chambers of Xeric: CM'),
    ('toa_solo_400', 'Tombs of Amascut'),
    ('toa_team_400', 'Tombs of Amascut')
) AS boss_list(boss, category);

DELETE FROM "bosses"
WHERE "name" IN ('maggot_king', 'mad_angel');

-- +goose StatementEnd
