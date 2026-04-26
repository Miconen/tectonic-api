-- +goose Up
-- +goose StatementBegin

-- Rename the table
ALTER TABLE "times" RENAME TO "records";

-- Rename columns
ALTER TABLE "records" RENAME COLUMN "time" TO "value";
ALTER TABLE "records" RENAME COLUMN "run_id" TO "record_id";

-- Rename the sequence
ALTER SEQUENCE "times_run_id_seq" RENAME TO "records_record_id_seq";

-- Rename constraints on records table
ALTER TABLE "records" RENAME CONSTRAINT "times_pkey" TO "records_pkey";
ALTER TABLE "records" RENAME CONSTRAINT "times_bosses_name_fkey" TO "records_boss_name_fkey";

-- Rename the column in teams table
ALTER TABLE "teams" RENAME COLUMN "run_id" TO "record_id";

-- Rename constraints on teams table
ALTER TABLE "teams" RENAME CONSTRAINT "teams_run_id_user_id_guild_id" TO "teams_record_id_user_id_guild_id";
ALTER TABLE "teams" RENAME CONSTRAINT "teams_run_id_fkey" TO "teams_record_id_fkey";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert teams constraints
ALTER TABLE "teams" RENAME CONSTRAINT "teams_record_id_fkey" TO "teams_run_id_fkey";
ALTER TABLE "teams" RENAME CONSTRAINT "teams_record_id_user_id_guild_id" TO "teams_run_id_user_id_guild_id";

-- Revert teams column
ALTER TABLE "teams" RENAME COLUMN "record_id" TO "run_id";

-- Revert records constraints
ALTER TABLE "records" RENAME CONSTRAINT "records_boss_name_fkey" TO "times_bosses_name_fkey";
ALTER TABLE "records" RENAME CONSTRAINT "records_pkey" TO "times_pkey";

-- Revert sequence
ALTER SEQUENCE "records_record_id_seq" RENAME TO "times_run_id_seq";

-- Revert columns
ALTER TABLE "records" RENAME COLUMN "record_id" TO "run_id";
ALTER TABLE "records" RENAME COLUMN "value" TO "time";

-- Revert table name
ALTER TABLE "records" RENAME TO "times";

-- +goose StatementEnd
