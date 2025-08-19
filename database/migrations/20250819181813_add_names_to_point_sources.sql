-- +goose Up
-- +goose StatementBegin
ALTER TABLE "point_sources"
ADD COLUMN "name" character varying(64);

UPDATE "point_sources"
SET "name" = point_sources.source
WHERE "name" IS NULL;

ALTER TABLE "point_sources"
ALTER COLUMN "name" SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "point_sources"
DROP COLUMN "name";
-- +goose StatementEnd
