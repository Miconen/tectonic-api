-- +goose Up
-- +goose StatementBegin
ALTER TABLE "event_participant"
ADD CONSTRAINT "event_participant_event_id_user_id_guild_id" PRIMARY KEY ("event_id", "user_id", "guild_id"),
DROP CONSTRAINT "event_participant_pkey";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "event_participant"
ADD CONSTRAINT "event_participant_event_id_user_id" PRIMARY KEY ("event_id", "user_id"),
DROP CONSTRAINT "event_participant_event_id_user_id_guild_id";
-- +goose StatementEnd
