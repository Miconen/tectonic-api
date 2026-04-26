-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."value_types" (
    "name" character varying(32) NOT NULL,
    "higher_is_better" boolean NOT NULL DEFAULT false,
    CONSTRAINT "value_types_pkey" PRIMARY KEY ("name")
) WITH (oids = false);

INSERT INTO "value_types" ("name", "higher_is_better")
VALUES
    ('time', false),
    ('depth', true);

ALTER TABLE "bosses"
ADD COLUMN "value_type" character varying(32) NOT NULL DEFAULT 'time';

ALTER TABLE "bosses"
ADD CONSTRAINT "bosses_value_type_fkey" FOREIGN KEY ("value_type")
REFERENCES "value_types" ("name") ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "bosses" DROP CONSTRAINT IF EXISTS "bosses_value_type_fkey";
ALTER TABLE "bosses" DROP COLUMN IF EXISTS "value_type";
DROP TABLE IF EXISTS "value_types";
-- +goose StatementEnd
