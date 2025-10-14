-- +goose Up
-- +goose StatementBegin
INSERT INTO achievement ("name", "thumbnail", "discord_icon", "order")
VALUES
    ('Gilded Collection Logger', 'https://oldschool.runescape.wiki/images/Collection_log_(gilded)_detail.png','<:glog:1427763493980864624>', 4);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM achievement
WHERE name = 'Gilded Collection Logger';
-- +goose StatementEnd
