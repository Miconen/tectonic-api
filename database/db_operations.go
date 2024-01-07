package database

import (
	"context"
	"fmt"
	"os"

	"tectonic-api/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
var db *pgx.Conn

func InitDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	db = conn // Store the connection in a package-level variable

	return conn, nil
}

func FetchUser(gid string, uid string) (models.User, error) {
	// Construct the SELECT query using Squirrel
	sql, args, err := psql.Select("*").From("users").Where(squirrel.Eq{"guild_id": gid, "user_id": uid}).ToSql()

	// Executing the query and fetching rows
	row := db.QueryRow(context.Background(), sql, args...)

	var user models.User

	// Scan data into User struct
	err = row.Scan(&user.UserId, &user.GuildId, &user.Points)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func FetchUsers(gid string) ([]models.User, error) {
	// Construct the SELECT query using Squirrel
	sql, args, err := psql.Select("*").From("users").Where(squirrel.Eq{"guild_id": gid}).ToSql()

	// Executing the query and fetching rows
	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	// Iterate through the rows and scan data into User struct
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserId, &user.GuildId, &user.Points); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func InsertUser(gid string, uid string) error {
	// Construct the SELECT query using Squirrel
	sql, args, err := psql.Insert("users").Columns("user_id", "guild_id").Values(uid, gid).ToSql()
	if err != nil {
		return err
	}

	// Executing the query using conn.Exec() for INSERT operation
	commandTag, err := db.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	// Check the command tag to verify the success of the INSERT operation
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", commandTag.RowsAffected())
	}

	return nil
}
