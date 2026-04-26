-- +goose Up
-- +goose StatementBegin
ALTER TABLE "guilds"
ADD COLUMN "position_count" smallint NOT NULL DEFAULT 3;

ALTER TABLE "guilds"
ADD CONSTRAINT "guilds_position_count_check" CHECK ("position_count" BETWEEN 1 AND 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "guilds" DROP CONSTRAINT IF EXISTS "guilds_position_count_check";
ALTER TABLE "guilds" DROP COLUMN IF EXISTS "position_count";
-- +goose StatementEnd
