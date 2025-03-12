// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkPb = `-- name: CheckPb :one
SELECT gb.boss, t.time
FROM guild_bosses gb
LEFT JOIN times t ON gb.pb_id = t.run_id AND gb.guild_id = t.guild_id
WHERE gb.boss = $1
AND gb.guild_id = $2
`

type CheckPbParams struct {
	Boss    string `json:"boss"`
	GuildID string `json:"guild_id"`
}

type CheckPbRow struct {
	Boss string      `json:"boss"`
	Time pgtype.Int4 `json:"time"`
}

func (q *Queries) CheckPb(ctx context.Context, arg CheckPbParams) (CheckPbRow, error) {
	row := q.db.QueryRow(ctx, checkPb, arg.Boss, arg.GuildID)
	var i CheckPbRow
	err := row.Scan(&i.Boss, &i.Time)
	return i, err
}

const createGuild = `-- name: CreateGuild :one
INSERT INTO guilds (
  guild_id
) VALUES (
  $1
)
RETURNING guild_id, multiplier, pb_channel_id
`

func (q *Queries) CreateGuild(ctx context.Context, guildID string) (Guild, error) {
	row := q.db.QueryRow(ctx, createGuild, guildID)
	var i Guild
	err := row.Scan(&i.GuildID, &i.Multiplier, &i.PbChannelID)
	return i, err
}

const createRsn = `-- name: CreateRsn :exec
INSERT INTO rsn (
    guild_id,
    user_id,
    rsn,
    wom_id
) VALUES (
    $1,
    $2,
    $3,
    $4
)
`

type CreateRsnParams struct {
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
	Rsn     string `json:"rsn"`
	WomID   string `json:"wom_id"`
}

func (q *Queries) CreateRsn(ctx context.Context, arg CreateRsnParams) error {
	_, err := q.db.Exec(ctx, createRsn,
		arg.GuildID,
		arg.UserID,
		arg.Rsn,
		arg.WomID,
	)
	return err
}

const createTeam = `-- name: CreateTeam :exec
INSERT INTO teams (run_id, user_id, guild_id)
SELECT $1, unnest($2::text[]), $3
`

type CreateTeamParams struct {
	RunID   int32    `json:"run_id"`
	UserIds []string `json:"user_ids"`
	GuildID string   `json:"guild_id"`
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) error {
	_, err := q.db.Exec(ctx, createTeam, arg.RunID, arg.UserIds, arg.GuildID)
	return err
}

const createTime = `-- name: CreateTime :one
INSERT INTO times (
    time,
    boss_name,
    date,
    guild_id
)
VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING run_id
`

type CreateTimeParams struct {
	Time     int32            `json:"time"`
	BossName string           `json:"boss_name"`
	Date     pgtype.Timestamp `json:"date"`
	GuildID  string           `json:"guild_id"`
}

func (q *Queries) CreateTime(ctx context.Context, arg CreateTimeParams) (int32, error) {
	row := q.db.QueryRow(ctx, createTime,
		arg.Time,
		arg.BossName,
		arg.Date,
		arg.GuildID,
	)
	var run_id int32
	err := row.Scan(&run_id)
	return run_id, err
}

const createUser = `-- name: CreateUser :one
WITH inserted_users AS (
  INSERT INTO users (guild_id, user_id)
  VALUES ($1, $2)
  RETURNING guild_id, user_id, points
),
inserted_rsn AS (
  INSERT INTO rsn (guild_id, user_id, rsn, wom_id)
  VALUES ($1, $2, $3, $4)
  RETURNING guild_id, user_id, rsn, wom_id
)
SELECT
    u.guild_id,
    u.user_id,
    u.points,
    r.rsn,
    r.wom_id
FROM inserted_users u
JOIN inserted_rsn r
ON u.guild_id = r.guild_id AND u.user_id = r.user_id
`

type CreateUserParams struct {
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
	Rsn     string `json:"rsn"`
	WomID   string `json:"wom_id"`
}

type CreateUserRow struct {
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
	Points  int32  `json:"points"`
	Rsn     string `json:"rsn"`
	WomID   string `json:"wom_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.GuildID,
		arg.UserID,
		arg.Rsn,
		arg.WomID,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.GuildID,
		&i.UserID,
		&i.Points,
		&i.Rsn,
		&i.WomID,
	)
	return i, err
}

const deleteGuild = `-- name: DeleteGuild :one
DELETE FROM guilds
WHERE guild_id = $1 RETURNING guild_id, multiplier, pb_channel_id
`

func (q *Queries) DeleteGuild(ctx context.Context, guildID string) (Guild, error) {
	row := q.db.QueryRow(ctx, deleteGuild, guildID)
	var i Guild
	err := row.Scan(&i.GuildID, &i.Multiplier, &i.PbChannelID)
	return i, err
}

const deleteRsn = `-- name: DeleteRsn :execrows
DELETE FROM rsn r
WHERE r.guild_id = $1 AND r.user_id = $2 AND r.rsn = $3
`

type DeleteRsnParams struct {
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
	Rsn     string `json:"rsn"`
}

func (q *Queries) DeleteRsn(ctx context.Context, arg DeleteRsnParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteRsn, arg.GuildID, arg.UserID, arg.Rsn)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteTime = `-- name: DeleteTime :execrows
DELETE FROM times t
WHERE t.guild_id = $1 AND t.run_id = $2
`

type DeleteTimeParams struct {
	GuildID string `json:"guild_id"`
	RunID   int32  `json:"run_id"`
}

func (q *Queries) DeleteTime(ctx context.Context, arg DeleteTimeParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteTime, arg.GuildID, arg.RunID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteUserById = `-- name: DeleteUserById :execrows
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id = $2
`

type DeleteUserByIdParams struct {
	GuildID string `json:"guild_id"`
	UserID  string `json:"user_id"`
}

func (q *Queries) DeleteUserById(ctx context.Context, arg DeleteUserByIdParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteUserById, arg.GuildID, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteUserByRsn = `-- name: DeleteUserByRsn :execrows
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.rsn = $2
)
`

type DeleteUserByRsnParams struct {
	GuildID string `json:"guild_id"`
	Rsn     string `json:"rsn"`
}

func (q *Queries) DeleteUserByRsn(ctx context.Context, arg DeleteUserByRsnParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteUserByRsn, arg.GuildID, arg.Rsn)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteUserByWom = `-- name: DeleteUserByWom :execrows
DELETE FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.wom_id = $2
)
`

type DeleteUserByWomParams struct {
	GuildID string `json:"guild_id"`
	WomID   string `json:"wom_id"`
}

func (q *Queries) DeleteUserByWom(ctx context.Context, arg DeleteUserByWomParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteUserByWom, arg.GuildID, arg.WomID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getBosses = `-- name: GetBosses :many
SELECT name, display_name, category, solo FROM bosses
`

func (q *Queries) GetBosses(ctx context.Context) ([]Boss, error) {
	rows, err := q.db.Query(ctx, getBosses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Boss
	for rows.Next() {
		var i Boss
		if err := rows.Scan(
			&i.Name,
			&i.DisplayName,
			&i.Category,
			&i.Solo,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategories = `-- name: GetCategories :many
SELECT "thumbnail", "order", "name" FROM categories
`

func (q *Queries) GetCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.Query(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.Thumbnail, &i.Order, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDetailedGuild = `-- name: GetDetailedGuild :one
SELECT
    g.guild_id,
    g.pb_channel_id,
    (SELECT json_agg(tm) FROM teams tm 
     JOIN times t ON tm.run_id = t.run_id 
     WHERE t.guild_id = g.guild_id) AS teammates,
    (SELECT json_agg(t) FROM times t 
     WHERE t.guild_id = g.guild_id) AS pbs,
    (SELECT json_agg(b) FROM bosses b 
     JOIN guild_bosses gb ON b.name = gb.boss 
     WHERE gb.guild_id = g.guild_id) AS bosses,
    (SELECT json_agg(c) FROM categories c 
     JOIN guild_categories gc ON c.name = gc.category 
     WHERE gc.guild_id = g.guild_id) AS categories,
    (SELECT json_agg(gb) FROM guild_bosses gb 
     WHERE gb.guild_id = g.guild_id) AS guild_bosses,
    (SELECT json_agg(gc) FROM guild_categories gc 
     WHERE gc.guild_id = g.guild_id) AS guild_categories
FROM guilds g
WHERE g.guild_id = $1
`

type GetDetailedGuildRow struct {
	GuildID         string      `json:"guild_id"`
	PbChannelID     pgtype.Text `json:"pb_channel_id"`
	Teammates       []byte      `json:"teammates"`
	Pbs             []byte      `json:"pbs"`
	Bosses          []byte      `json:"bosses"`
	Categories      []byte      `json:"categories"`
	GuildBosses     []byte      `json:"guild_bosses"`
	GuildCategories []byte      `json:"guild_categories"`
}

func (q *Queries) GetDetailedGuild(ctx context.Context, guildID string) (GetDetailedGuildRow, error) {
	row := q.db.QueryRow(ctx, getDetailedGuild, guildID)
	var i GetDetailedGuildRow
	err := row.Scan(
		&i.GuildID,
		&i.PbChannelID,
		&i.Teammates,
		&i.Pbs,
		&i.Bosses,
		&i.Categories,
		&i.GuildBosses,
		&i.GuildCategories,
	)
	return i, err
}

const getDetailedUsers = `-- name: GetDetailedUsers :many
SELECT 
    du.user_id,
    du.guild_id,
    du.points,
    to_json(du.rsns) AS rsns,
    COALESCE(times_json, '[]'::json) AS times
FROM detailed_users du
LEFT JOIN LATERAL (
    SELECT json_agg(dt) AS times_json
    FROM detailed_times dt
    WHERE dt.run_id IN (
        SELECT tm.run_id
        FROM teams tm
        WHERE tm.user_id = du.user_id AND tm.guild_id = du.guild_id
    )
) t ON true
WHERE du.user_id = ANY($1::text[])
AND du.guild_id = $2
`

type GetDetailedUsersParams struct {
	UserIds []string `json:"user_ids"`
	GuildID string   `json:"guild_id"`
}

type GetDetailedUsersRow struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Points  int32  `json:"points"`
	Rsns    []byte `json:"rsns"`
	Times   []byte `json:"times"`
}

func (q *Queries) GetDetailedUsers(ctx context.Context, arg GetDetailedUsersParams) ([]GetDetailedUsersRow, error) {
	rows, err := q.db.Query(ctx, getDetailedUsers, arg.UserIds, arg.GuildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDetailedUsersRow
	for rows.Next() {
		var i GetDetailedUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.GuildID,
			&i.Points,
			&i.Rsns,
			&i.Times,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDetailedUsersByRSN = `-- name: GetDetailedUsersByRSN :many
WITH rsn_user AS (
    SELECT r.user_id, r.guild_id 
    FROM rsn r 
    WHERE r.rsn = ANY($2::text[])
)
SELECT 
    du.user_id,
    du.guild_id,
    du.points,
    to_json(du.rsns) AS rsns,
    COALESCE(times_json, '[]'::json) AS times
FROM detailed_users du
JOIN rsn_user ru ON ru.user_id = du.user_id AND ru.guild_id = du.guild_id
LEFT JOIN LATERAL (
    SELECT json_agg(dt) AS times_json
    FROM detailed_times dt
    WHERE dt.run_id IN (
        SELECT tm.run_id
        FROM teams tm
        WHERE tm.user_id = du.user_id AND tm.guild_id = du.guild_id
    )
) t ON true
WHERE du.guild_id = $1
`

type GetDetailedUsersByRSNParams struct {
	GuildID string   `json:"guild_id"`
	Rsns    []string `json:"rsns"`
}

type GetDetailedUsersByRSNRow struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Points  int32  `json:"points"`
	Rsns    []byte `json:"rsns"`
	Times   []byte `json:"times"`
}

func (q *Queries) GetDetailedUsersByRSN(ctx context.Context, arg GetDetailedUsersByRSNParams) ([]GetDetailedUsersByRSNRow, error) {
	rows, err := q.db.Query(ctx, getDetailedUsersByRSN, arg.GuildID, arg.Rsns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDetailedUsersByRSNRow
	for rows.Next() {
		var i GetDetailedUsersByRSNRow
		if err := rows.Scan(
			&i.UserID,
			&i.GuildID,
			&i.Points,
			&i.Rsns,
			&i.Times,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDetailedUsersByWomID = `-- name: GetDetailedUsersByWomID :many
WITH wom_user AS (
    SELECT r.user_id, r.guild_id 
    FROM rsn r 
    WHERE r.wom_id = ANY($2::text[])
)
SELECT 
    du.user_id,
    du.guild_id,
    du.points,
    to_json(du.rsns) AS rsns,
    COALESCE(times_json, '[]'::json) AS times
FROM detailed_users du
JOIN wom_user wu ON wu.user_id = du.user_id AND wu.guild_id = du.guild_id
LEFT JOIN LATERAL (
    SELECT json_agg(dt) AS times_json
    FROM detailed_times dt
    WHERE dt.run_id IN (
        SELECT tm.run_id
        FROM teams tm
        WHERE tm.user_id = du.user_id AND tm.guild_id = du.guild_id
    )
) t ON true
WHERE du.guild_id = $1
`

type GetDetailedUsersByWomIDParams struct {
	GuildID string   `json:"guild_id"`
	WomIds  []string `json:"wom_ids"`
}

type GetDetailedUsersByWomIDRow struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Points  int32  `json:"points"`
	Rsns    []byte `json:"rsns"`
	Times   []byte `json:"times"`
}

func (q *Queries) GetDetailedUsersByWomID(ctx context.Context, arg GetDetailedUsersByWomIDParams) ([]GetDetailedUsersByWomIDRow, error) {
	rows, err := q.db.Query(ctx, getDetailedUsersByWomID, arg.GuildID, arg.WomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDetailedUsersByWomIDRow
	for rows.Next() {
		var i GetDetailedUsersByWomIDRow
		if err := rows.Scan(
			&i.UserID,
			&i.GuildID,
			&i.Points,
			&i.Rsns,
			&i.Times,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGuild = `-- name: GetGuild :one
SELECT guild_id, multiplier, pb_channel_id FROM guilds
WHERE guild_id = $1 LIMIT 1
`

func (q *Queries) GetGuild(ctx context.Context, guildID string) (Guild, error) {
	row := q.db.QueryRow(ctx, getGuild, guildID)
	var i Guild
	err := row.Scan(&i.GuildID, &i.Multiplier, &i.PbChannelID)
	return i, err
}

const getLeaderboard = `-- name: GetLeaderboard :many
SELECT u.user_id, u.guild_id, u.points, json_agg(r) AS rsns
FROM users u
JOIN rsn r ON u.user_id = r.user_id AND u.guild_id = r.guild_id
WHERE u.guild_id = $1
GROUP BY u.user_id, u.guild_id, u.points
ORDER BY u.points DESC
LIMIT $2
`

type GetLeaderboardParams struct {
	GuildID   string `json:"guild_id"`
	UserLimit int32  `json:"user_limit"`
}

type GetLeaderboardRow struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
	Points  int32  `json:"points"`
	Rsns    []byte `json:"rsns"`
}

func (q *Queries) GetLeaderboard(ctx context.Context, arg GetLeaderboardParams) ([]GetLeaderboardRow, error) {
	rows, err := q.db.Query(ctx, getLeaderboard, arg.GuildID, arg.UserLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLeaderboardRow
	for rows.Next() {
		var i GetLeaderboardRow
		if err := rows.Scan(
			&i.UserID,
			&i.GuildID,
			&i.Points,
			&i.Rsns,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPointsValue = `-- name: GetPointsValue :one
SELECT points
FROM point_sources
WHERE source = $1
AND guild_id = $2
`

type GetPointsValueParams struct {
	Event   string `json:"event"`
	GuildID string `json:"guild_id"`
}

func (q *Queries) GetPointsValue(ctx context.Context, arg GetPointsValueParams) (int32, error) {
	row := q.db.QueryRow(ctx, getPointsValue, arg.Event, arg.GuildID)
	var points int32
	err := row.Scan(&points)
	return points, err
}

const getUsersById = `-- name: GetUsersById :many
SELECT users.user_id, users.guild_id, users.points
FROM users
WHERE users.guild_id = $1
AND users.user_id = ANY($2::text[])
`

type GetUsersByIdParams struct {
	GuildID string   `json:"guild_id"`
	UserIds []string `json:"user_ids"`
}

func (q *Queries) GetUsersById(ctx context.Context, arg GetUsersByIdParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersById, arg.GuildID, arg.UserIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.UserID, &i.GuildID, &i.Points); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByRsn = `-- name: GetUsersByRsn :many
SELECT users.user_id, users.guild_id, users.points
FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.rsn = ANY($2::text[])
)
`

type GetUsersByRsnParams struct {
	GuildID string   `json:"guild_id"`
	Rsns    []string `json:"rsns"`
}

func (q *Queries) GetUsersByRsn(ctx context.Context, arg GetUsersByRsnParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersByRsn, arg.GuildID, arg.Rsns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.UserID, &i.GuildID, &i.Points); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByWom = `-- name: GetUsersByWom :many
SELECT users.user_id, users.guild_id, users.points
FROM users
WHERE users.guild_id = $1
AND users.user_id IN (
    SELECT rsn.user_id
    FROM rsn
    WHERE rsn.guild_id = users.guild_id AND rsn.wom_id = ANY($2::text[])
)
`

type GetUsersByWomParams struct {
	GuildID string   `json:"guild_id"`
	WomIds  []string `json:"wom_ids"`
}

func (q *Queries) GetUsersByWom(ctx context.Context, arg GetUsersByWomParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersByWom, arg.GuildID, arg.WomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.UserID, &i.GuildID, &i.Points); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCategoryMessageIds = `-- name: UpdateCategoryMessageIds :execrows
UPDATE guild_categories
SET message_id = u.message_id
FROM (SELECT unnest($2::text[]) as category,
             unnest($3::text[]) as message_id) as u
WHERE guild_categories.guild_id = $1
AND guild_categories.category = u.category
`

type UpdateCategoryMessageIdsParams struct {
	GuildID    string   `json:"guild_id"`
	Categories []string `json:"categories"`
	MessageIds []string `json:"message_ids"`
}

func (q *Queries) UpdateCategoryMessageIds(ctx context.Context, arg UpdateCategoryMessageIdsParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateCategoryMessageIds, arg.GuildID, arg.Categories, arg.MessageIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updateGuild = `-- name: UpdateGuild :one
UPDATE guilds SET
    multiplier = CASE WHEN $1::numeric IS NOT NULL AND $1::numeric != 0 THEN $1::numeric ELSE multiplier END,
    pb_channel_id = CASE WHEN $2::text IS NOT NULL AND $2::text != '' THEN $2::text ELSE pb_channel_id END
WHERE guild_id = $3 RETURNING guild_id, multiplier, pb_channel_id
`

type UpdateGuildParams struct {
	Multiplier  pgtype.Numeric `json:"multiplier"`
	PbChannelID string         `json:"pb_channel_id"`
	GuildID     string         `json:"guild_id"`
}

func (q *Queries) UpdateGuild(ctx context.Context, arg UpdateGuildParams) (Guild, error) {
	row := q.db.QueryRow(ctx, updateGuild, arg.Multiplier, arg.PbChannelID, arg.GuildID)
	var i Guild
	err := row.Scan(&i.GuildID, &i.Multiplier, &i.PbChannelID)
	return i, err
}

const updatePb = `-- name: UpdatePb :execrows
UPDATE guild_bosses SET
    pb_id = $1
WHERE guild_id = $2
AND boss = $3
`

type UpdatePbParams struct {
	RunID   pgtype.Int4 `json:"run_id"`
	GuildID string      `json:"guild_id"`
	Boss    string      `json:"boss"`
}

func (q *Queries) UpdatePb(ctx context.Context, arg UpdatePbParams) (int64, error) {
	result, err := q.db.Exec(ctx, updatePb, arg.RunID, arg.GuildID, arg.Boss)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updatePointsByEvent = `-- name: UpdatePointsByEvent :many
WITH point_value AS (
    SELECT points
    FROM point_sources
    WHERE source = $3
    AND guild_id = $2
)
UPDATE users
SET points = points + (SELECT points FROM point_value)
WHERE user_id = ANY($1::text[])
AND users.guild_id = $2 
RETURNING user_id, guild_id, points, (SELECT points FROM point_value) AS given_points
`

type UpdatePointsByEventParams struct {
	UserIds []string `json:"user_ids"`
	GuildID string   `json:"guild_id"`
	Event   string   `json:"event"`
}

type UpdatePointsByEventRow struct {
	UserID      string `json:"user_id"`
	GuildID     string `json:"guild_id"`
	Points      int32  `json:"points"`
	GivenPoints int32  `json:"given_points"`
}

func (q *Queries) UpdatePointsByEvent(ctx context.Context, arg UpdatePointsByEventParams) ([]UpdatePointsByEventRow, error) {
	rows, err := q.db.Query(ctx, updatePointsByEvent, arg.UserIds, arg.GuildID, arg.Event)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UpdatePointsByEventRow
	for rows.Next() {
		var i UpdatePointsByEventRow
		if err := rows.Scan(
			&i.UserID,
			&i.GuildID,
			&i.Points,
			&i.GivenPoints,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePointsCustom = `-- name: UpdatePointsCustom :execrows
UPDATE users
SET points = points + $1
WHERE user_id = ANY($2::text[])
AND guild_id = $3 RETURNING user_id, guild_id, points
`

type UpdatePointsCustomParams struct {
	Points  int32    `json:"points"`
	UserIds []string `json:"user_ids"`
	GuildID string   `json:"guild_id"`
}

func (q *Queries) UpdatePointsCustom(ctx context.Context, arg UpdatePointsCustomParams) (int64, error) {
	result, err := q.db.Exec(ctx, updatePointsCustom, arg.Points, arg.UserIds, arg.GuildID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
