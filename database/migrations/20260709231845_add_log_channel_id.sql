-- +goose Up
-- +goose StatementBegin
ALTER TABLE "guilds"
ADD "log_channel_id" character varying(32);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "guilds"
DROP COLUMN "log_channel_id";
-- +goose StatementEnd
