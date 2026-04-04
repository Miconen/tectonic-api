-- +goose Up
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('demonic_brutus', 'Demonic Brutus', 'Miscellaneous', true);

INSERT INTO "guild_bosses" ("boss", "guild_id", "pb_id", "category")
SELECT 'demonic_brutus', g.guild_id, NULL, 'Miscellaneous'
FROM "guilds" g;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses" WHERE "name" = 'demonic_brutus';
-- +goose StatementEnd
