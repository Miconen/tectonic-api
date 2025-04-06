-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."achievement" (
    "name" character varying(32) NOT NULL,
    "thumbnail" character varying(256) NOT NULL,
    "discord_icon" character varying(32) NOT NULL,
    "order" smallint DEFAULT '0' NOT NULL,
    CONSTRAINT "achievement_name" PRIMARY KEY ("name")
) WITH (oids = false);

INSERT INTO achievement ("name", "thumbnail", "discord_icon", "order")
VALUES
    -- Ironman status
    ('Ironman', 'https://oldschool.runescape.wiki/images/Ironman_chat_badge.png','<:IM:990015824981020702>', 1),
    ('HCIM', 'https://oldschool.runescape.wiki/images/Hardcore_ironman_chat_badge.png','<:HCIM:990015822376366140>', 1),
    ('UIM', 'https://oldschool.runescape.wiki/images/Ultimate_ironman_chat_badge.png','<:UIM:990015823651422238>', 1),
    ('GIM', 'https://oldschool.runescape.wiki/images/Group_ironman_chat_badge.png','<:GIM:990015820568592434>', 1),
    ('HCGIM', 'https://oldschool.runescape.wiki/images/Hardcore_group_ironman_chat_badge.png','<:HCGIM:990015818924429312>', 1),

    -- Actual achievements
    ('Maxed', 'https://oldschool.runescape.wiki/images/Max_cape.png','<:MaxCape:1332163641071505459>', 2),
    ('Grandmaster', 'https://oldschool.runescape.wiki/images/Tzkal_slayer_helmet.png','<:ZukHelm:979469356289392640>', 3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "achievement";
-- +goose StatementEnd
