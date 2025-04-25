-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user_achievement"
ADD "guild_id" character varying(32) NOT NULL;

ALTER TABLE "user_achievement"
ADD CONSTRAINT "user_achievement_user_id_achievement_name_guild_id" PRIMARY KEY ("user_id", "achievement_name", "guild_id"),
DROP CONSTRAINT "user_achievement_pkey";

ALTER TABLE "user_achievement" ADD FOREIGN KEY ("user_id", "guild_id") REFERENCES "users" ("user_id", "guild_id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_achievement" ADD FOREIGN KEY ("achievement_name") REFERENCES "achievement" ("name") ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user_achievement" DROP CONSTRAINT "user_achievement_user_id_guild_id_fkey";
ALTER TABLE "user_achievement" DROP CONSTRAINT "user_achievement_achievement_name_fkey";

ALTER TABLE "user_achievement"
ADD CONSTRAINT "user_achievement_user_id_achievement_name" PRIMARY KEY ("user_id", "achievement_name"),
DROP CONSTRAINT "user_achievement_user_id_achievement_name_guild_id";

ALTER TABLE "user_achievement"
DROP "guild_id";
-- +goose StatementEnd
