-- +goose Up
-- +goose StatementBegin
INSERT INTO "bosses" ("name", "display_name", "category", "solo")
VALUES ('araxxor', 'Araxxor', 'Slayer Boss', '1');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "bosses"
WHERE "name" = 'araxxor';
-- +goose StatementEnd
