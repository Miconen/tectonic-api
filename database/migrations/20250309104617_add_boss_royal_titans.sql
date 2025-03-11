-- +goose Up
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('royal_titans_1', 'Royal Titans (Solo)', 'Miscellaneous', '1');

INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('royal_titans_2', 'Royal Titans (Duo)', 'Miscellaneous', '0');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses"
WHERE (("name" = 'royal_titans_1') OR ("name" = 'royal_titans_2'));
-- +goose StatementEnd
