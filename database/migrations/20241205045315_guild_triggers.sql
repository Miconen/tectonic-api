-- +goose Up
-- +goose StatementBegin
ALTER TABLE guild_categories
ALTER COLUMN message_id DROP NOT NULL;

UPDATE guild_categories
SET message_id = NULL
WHERE message_id = '';

-- Function to insert guild_bosses and guild_categories for each new guild
CREATE OR REPLACE FUNCTION insert_guild_bosses_and_categories()
RETURNS TRIGGER AS $$
BEGIN
  -- Insert selected rows from categories table into guild_categories
  INSERT INTO guild_categories (guild_id, category)
  SELECT NEW.guild_id, name
  FROM categories;

  -- Insert selected rows from bosses table into guild_bosses
  INSERT INTO guild_bosses (guild_id, boss)
  SELECT NEW.guild_id, name
  FROM bosses;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to invoke the function after a new guild is inserted
CREATE TRIGGER insert_guild_bosses_and_categories_trigger
AFTER INSERT ON guilds
FOR EACH ROW
EXECUTE FUNCTION insert_guild_bosses_and_categories();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE guild_categories
SET message_id = ''
WHERE message_id = NULL;

ALTER TABLE guild_categories
ALTER COLUMN message_id SET NOT NULL;

DROP TRIGGER insert_guild_bosses_and_categories_trigger ON guilds;
DROP FUNCTION insert_guild_bosses_and_categories();
-- +goose StatementEnd
