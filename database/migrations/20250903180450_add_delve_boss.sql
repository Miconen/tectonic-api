-- +goose Up
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('doom_of_mokhaiotl', 'Doom of Mokhaiotl (1-8)', 'Varlamore', '1');

INSERT INTO "guild_bosses" ("boss", "guild_id", "pb_id")
SELECT 'doom_of_mokhaiotl', g.guild_id, NULL
FROM "guilds" g;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses" WHERE "name" = 'doom_of_mokhaiotl';
-- +goose StatementEnd
