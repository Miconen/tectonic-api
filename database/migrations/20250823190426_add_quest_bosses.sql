-- +goose Up
-- +goose StatementBegin
INSERT INTO "categories" ("thumbnail", "order", "name")
VALUES ('https://oldschool.runescape.wiki/images/Quests.png', 13, 'Quest Bosses');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('seren', 'Seren', 'Quest Bosses', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('galvek', 'Galvek', 'Quest Bosses', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('glough', 'Glough', 'Quest Bosses', '1');

INSERT INTO "guild_categories" ("guild_id", "category", "message_id")
SELECT g.guild_id, category_list.category, NULL
FROM "guilds" g
CROSS JOIN (VALUES ('Quest Bosses')) AS category_list(category);

INSERT INTO "guild_bosses" ("boss", "guild_id", "pb_id")
SELECT boss_list.boss, g.guild_id, NULL
FROM "guilds" g
CROSS JOIN (VALUES ('seren'), ('galvek'), ('glough')) AS boss_list(boss);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses"
WHERE (("name" = 'seren') OR ("name" = 'galvek') OR ("name" = 'glough'));

DELETE FROM "categories"
WHERE "cateogry" = 'Quest Bosses';
-- +goose StatementEnd
