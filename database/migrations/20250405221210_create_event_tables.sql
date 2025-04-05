-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."event" (
    "name" character varying(64) NOT NULL,
    "wom_id" character varying(32) NOT NULL,
    "guild_id" character varying(32) NOT NULL,
    "position_cutoff" smallint DEFAULT '1' NOT NULL,
    CONSTRAINT "event_pkey" PRIMARY KEY ("wom_id")
) WITH (oids = false);

CREATE TABLE "public"."event_participant" (
    "user_id" character varying(32) NOT NULL,
    "guild_id" character varying(32) NOT NULL,
    "placement" smallint NOT NULL,
    "event_id" character varying(32) NOT NULL,
    CONSTRAINT "event_participant_pkey" PRIMARY KEY ("event_id", "user_id")
) WITH (oids = false);

ALTER TABLE ONLY "public"."event_participant" ADD CONSTRAINT "event_id_fkey" FOREIGN KEY (event_id) REFERENCES event(wom_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
ALTER TABLE ONLY "public"."event" ADD CONSTRAINT "event_guild_id_fkey" FOREIGN KEY (guild_id) REFERENCES guilds(guild_id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "event_participant";
DROP TABLE IF EXISTS "event";
-- +goose StatementEnd
