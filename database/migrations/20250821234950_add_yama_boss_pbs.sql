-- +goose Up
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('yama_1', 'Yama (Solo)', 'Miscellaneous', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('yama_2', 'Yama (Duo)', 'Miscellaneous', '0');

INSERT INTO "guild_bosses" ("boss", "guild_id", "pb_id")
SELECT boss_list.boss, g.guild_id, NULL
FROM "guilds" g
CROSS JOIN (VALUES ('yama_1'), ('yama_2')) AS boss_list(boss);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses"
WHERE (("name" = 'yama_1') OR ("name" = 'yama_2'));
-- +goose StatementEnd
