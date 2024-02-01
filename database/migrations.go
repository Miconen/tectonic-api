package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/z0ne-dev/mgx"
)

func RunMigrations(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	addGamemodeColumnToRsnTable := mgx.NewRawMigration("addGamemodeColumnToRsnTable",
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gamemode') THEN
				CREATE TYPE gamemode AS ENUM ('main', 'ironman', 'ultimate_ironman', 'hardcore_ironman', 'group_ironman', 'unranked_group_ironman', 'hardcore_group_ironman');
			END IF;
		END
		$$;

		ALTER TABLE public.rsn
		DROP COLUMN IF EXISTS gamemode;

		ALTER TABLE public.rsn
		ADD COLUMN gamemode gamemode DEFAULT 'main' NOT NULL;`)

	migrator, err := mgx.New(mgx.Migrations(addGamemodeColumnToRsnTable))
	if err != nil {
		return err
	}

	err = migrator.Migrate(context.TODO(), conn)
	if err != nil {
		return err
	}

	return nil
}
