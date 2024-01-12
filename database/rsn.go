package database

import (
	"context"
	"fmt"
	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
)

func SelectRsns(f map[string]string) ([]models.RSN, error) {
	query := psql.Select("*").From("rsn")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return []models.RSN{}, err
	}

	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return []models.RSN{}, err
	}

	var rsns []models.RSN

	// Iterate through the rows and scan data into User struct
	for rows.Next() {
		var rsn models.RSN
		if err := rows.Scan(&rsn.UserId, &rsn.GuildId, &rsn.WomId, &rsn.RSN); err != nil {
			return []models.RSN{}, err
		}
		rsns = append(rsns, rsn)
	}

	return rsns, nil
}

func InsertRsn(f map[string]string, wid string) error {
	query := psql.Insert("rsn").Columns("guild_id", "user_id", "rsn", "wom_id").Values(f["guild_id"], f["user_id"], f["rsn"], wid)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}

func DeleteRsn(f map[string]string) error {
	query := psql.Delete("rsn")

	for key, value := range f {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}
