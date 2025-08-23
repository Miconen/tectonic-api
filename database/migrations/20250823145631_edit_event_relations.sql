-- +goose Up
-- +goose StatementBegin
-- Drop the old foreign key constraint
ALTER TABLE event_participant
DROP CONSTRAINT event_id_fkey;

-- Drop the current composite primary key
ALTER TABLE event
DROP CONSTRAINT event_pkey;

-- Create new primary key
ALTER TABLE event
ADD CONSTRAINT event_pkey
PRIMARY KEY (wom_id, guild_id);

-- Create new foreign key constraint
ALTER TABLE event_participant
ADD CONSTRAINT event_id_guild_id_fkey
FOREIGN KEY (event_id, guild_id) 
REFERENCES event(wom_id, guild_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the new foreign key constraint
ALTER TABLE event_participant
DROP CONSTRAINT event_id_guild_id_fkey;

-- Drop the current composite primary key
ALTER TABLE event
DROP CONSTRAINT event_pkey;

-- Recreate the original primary key
ALTER TABLE event
ADD CONSTRAINT event_pkey
PRIMARY KEY (wom_id);

-- Recreate the original foreign key constraint
ALTER TABLE event_participant
ADD CONSTRAINT event_id_fkey
FOREIGN KEY (event_id) 
REFERENCES event(wom_id);
-- +goose StatementEnd
