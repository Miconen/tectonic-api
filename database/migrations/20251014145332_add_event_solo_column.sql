-- +goose Up
-- +goose StatementBegin
ALTER TABLE "event"
ADD "solo" boolean;

UPDATE "event" SET "solo" = TRUE;

ALTER TABLE "event"
ALTER COLUMN "solo" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "event"
DROP "solo";
-- +goose StatementEnd
