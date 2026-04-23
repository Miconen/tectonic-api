-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS insert_default_combat_achievements_trigger ON guilds;

CREATE TRIGGER seed_default_combat_achievements_trigger
AFTER INSERT ON guilds
FOR EACH ROW
EXECUTE FUNCTION insert_default_combat_achievements();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS seed_default_combat_achievements_trigger ON guilds;

CREATE TRIGGER insert_default_combat_achievements_trigger
AFTER INSERT ON guilds
FOR EACH ROW
EXECUTE FUNCTION insert_default_combat_achievements();
-- +goose StatementEnd
