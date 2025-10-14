-- +goose Up
-- +goose StatementBegin
ALTER TABLE "event_participant" DROP CONSTRAINT "event_id_guild_id_fkey";
ALTER TABLE "event_participant" ADD FOREIGN KEY ("event_id", "guild_id") REFERENCES "event" ("wom_id", "guild_id") ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "event_participant" DROP CONSTRAINT "event_participant_event_id_guild_id_fkey";
ALTER TABLE "event_participant" ADD FOREIGN KEY ("event_id", "guild_id") REFERENCES "event" ("wom_id", "guild_id");
-- +goose StatementEnd
