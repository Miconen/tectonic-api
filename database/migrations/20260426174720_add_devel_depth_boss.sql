-- +goose Up
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo", "value_type")
VALUES ('doom_of_mokhaiotl_depth', 'Doom of Mokhaiotl (Depth)', 'Varlamore', true, 'depth');

INSERT INTO "guild_bosses" ("boss", "guild_id", "category")
SELECT 'doom_of_mokhaiotl_depth', g.guild_id, 'Varlamore'
FROM "guilds" g;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses" WHERE "name" = 'doom_of_mokhaiotl_depth';
-- +goose StatementEnd
