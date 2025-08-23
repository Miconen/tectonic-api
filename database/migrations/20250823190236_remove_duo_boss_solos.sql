-- +goose Up
-- +goose StatementBegin
DELETE FROM "bosses"
WHERE (("name" = 'yama_1') OR ("name" = 'royal_titans_1'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('yama_1', 'Yama (Solo)', 'Miscellaneous', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('royal_titans_1', 'Royal Titans (Solo)', 'Miscellaneous', '1');

INSERT INTO "guild_bosses" ("boss", "guild_id", "pb_id")
SELECT boss_list.boss, g.guild_id, NULL
FROM "guilds" g
CROSS JOIN (VALUES ('yama_1'), ('royal_titans_1')) AS boss_list(boss);
-- +goose StatementEnd
