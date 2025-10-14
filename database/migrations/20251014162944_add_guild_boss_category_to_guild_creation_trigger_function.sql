-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION insert_guild_bosses_and_categories()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO guild_categories (guild_id, category)
  SELECT NEW.guild_id, name
  FROM categories;

  INSERT INTO guild_bosses (guild_id, boss, category)
  SELECT NEW.guild_id, name, category
  FROM bosses;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION insert_guild_bosses_and_categories()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO guild_categories (guild_id, category)
  SELECT NEW.guild_id, name
  FROM categories;

  INSERT INTO guild_bosses (guild_id, boss)
  SELECT NEW.guild_id, name
  FROM bosses;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
