-- +goose Up
-- +goose StatementBegin

-- Drop the FK constraint first, then the column
ALTER TABLE "guild_bosses" DROP CONSTRAINT IF EXISTS "guild_bosses_pb_id_fkey";
ALTER TABLE "guild_bosses" DROP COLUMN "pb_id";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Re-add pb_id column (nullable, references records which was formerly times)
ALTER TABLE "guild_bosses"
ADD COLUMN "pb_id" integer;

ALTER TABLE "guild_bosses"
ADD CONSTRAINT "guild_bosses_pb_id_fkey" FOREIGN KEY ("pb_id"
REFERENCES "records" ("record_id") ON UPDATE CASCADE ON DELETE SET NULL NOT DEFERRABLE;

-- +goose StatementEnd
