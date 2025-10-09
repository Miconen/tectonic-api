-- +goose Up
-- +goose StatementBegin
ALTER TABLE "guild_bosses"
ADD "category" character varying(64);

UPDATE "guild_bosses"
SET "category" = bosses.category
FROM bosses
WHERE guild_bosses.boss = bosses.name;

ALTER TABLE "guild_bosses"
ALTER COLUMN "category" SET NOT NULL;

ALTER TABLE "guild_bosses" ADD FOREIGN KEY ("category") REFERENCES "categories" ("name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "guild_bosses"
DROP "category";
-- +goose StatementEnd
