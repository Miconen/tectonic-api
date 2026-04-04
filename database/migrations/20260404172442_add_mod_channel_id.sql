-- +goose Up
-- +goose StatementBegin
ALTER TABLE "guilds"
ADD "mod_channel_id" character varying(32);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "guilds"
DROP COLUMN "mod_channel_id";
-- +goose StatementEnd
