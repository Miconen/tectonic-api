-- +goose Up
-- +goose StatementBegin
ALTER TABLE "times"
ADD "guild_id" character varying(32) NOT NULL;

ALTER TABLE "times"
ADD FOREIGN KEY ("guild_id") REFERENCES "guilds" ("guild_id") ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "times"
DROP "guild_id";

ALTER TABLE "times"
DROP CONSTRAINT "guild_id";
-- +goose StatementEnd
