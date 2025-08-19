-- +goose Up
-- +goose StatementBegin
DELETE FROM "achievement"
WHERE (("name" = 'Ironman') OR ("name" = 'HCIM') OR ("name" = 'UIM') OR ("name" = 'GIM') OR ("name" = 'HCGIM'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO achievement ("name", "thumbnail", "discord_icon", "order")
VALUES
    -- Ironman status
    ('Ironman', 'https://oldschool.runescape.wiki/images/Ironman_chat_badge.png','<:IM:990015824981020702>', 1),
    ('HCIM', 'https://oldschool.runescape.wiki/images/Hardcore_ironman_chat_badge.png','<:HCIM:990015822376366140>', 1),
    ('UIM', 'https://oldschool.runescape.wiki/images/Ultimate_ironman_chat_badge.png','<:UIM:990015823651422238>', 1),
    ('GIM', 'https://oldschool.runescape.wiki/images/Group_ironman_chat_badge.png','<:GIM:990015820568592434>', 1),
    ('HCGIM', 'https://oldschool.runescape.wiki/images/Hardcore_group_ironman_chat_badge.png','<:HCGIM:990015818924429312>', 1),
-- +goose StatementEnd
