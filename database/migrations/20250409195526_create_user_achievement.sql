-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."user_achievement" (
	"user_id" VARCHAR(32) NOT NULL,
	"achievement_name" VARCHAR(32) NOT NULL,
	CONSTRAINT "user_achievement_pkey" PRIMARY KEY ("user_id", "achievement_name")
) WITH (oids = false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_achievement;
-- +goose StatementEnd
