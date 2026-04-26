-- +goose Up
-- +goose StatementBegin
CREATE INDEX "idx_records_guild_boss_value" ON "records" ("guild_id", "boss_name", "value");
CREATE INDEX "idx_records_guild_boss_date" ON "records" ("guild_id", "boss_name", "date" DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_records_guild_boss_date";
DROP INDEX IF EXISTS "idx_records_guild_boss_value";
-- +goose StatementEnd
