package database

import (
	"context"
)

// TODO: Move to using transactions or something more sensible, just leaving this here for now
// BEGIN;
//
// SELECT * FROM users WHERE user_id = $1 AND guild_id = $2;
// SELECT * FROM rsn WHERE user_id = $1 AND guild_id = $2;
//
// SELECT
//     teams.run_id,
//     times.time,
//     times.boss_name,
//     times.date,
//     users.user_id,
//     users.guild_id,
//     users.points
// FROM teams
// JOIN times ON teams.run_id = times.run_id
// JOIN users ON teams.user_id = users.user_id
// WHERE times.run_id IN (
//     SELECT pb_id
//     FROM guild_bosses
//     WHERE guild_bosses.guild_id = $2
// );
//
// COMMIT;

const getUsers = `-- name: GetUsers :many
SELECT json_agg(
    json_build_object(
        'user', json_build_object(
            'user_id', u.user_id,
            'guild_id', u.guild_id,
            'points', u.points,
            'rsns', (
                SELECT json_agg(
                    json_build_object(
                        'rsn', r.rsn,
                        'wom_id', r.wom_id,
                        'user_id', r.user_id,
                        'guild_id', r.guild_id
                    )
                )
                FROM rsn r
                WHERE r.user_id = u.user_id AND r.guild_id = u.guild_id
            )
        ),
        'times', (
            SELECT json_agg(
                json_build_object(
                    'time', t.time,
                    'boss_name', t.boss_name,
                    'run_id', t.run_id,
                    'date', t.date,
                    'team', (
                        SELECT json_agg(
                            json_build_object(
                                'user_id', tu.user_id,
                                'guild_id', tu.guild_id,
                                'points', tu.points,
                                'rsns', (
                                    SELECT json_agg(
                                        json_build_object(
                                            'rsn', tr.rsn,
                                            'wom_id', tr.wom_id
                                        )
                                    )
                                    FROM rsn tr
                                    WHERE tr.user_id = tu.user_id AND tr.guild_id = tu.guild_id
                                )
                            )
                        )
                        FROM teams tm
                        JOIN users tu ON tm.user_id = tu.user_id AND tm.guild_id = tu.guild_id
                        WHERE tm.run_id = t.run_id
                    )
                )
            )
            FROM times t
            WHERE t.run_id IN (
                SELECT gb.pb_id
                FROM guild_bosses gb
                WHERE gb.guild_id = u.guild_id
            )
        )
    )
)
AS result
FROM users u;
`

func (q *Queries) GetUsers(ctx context.Context) ([][]byte, error) {
	rows, err := q.db.Query(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]byte
	for rows.Next() {
		var result []byte
		if err := rows.Scan(&result); err != nil {
			return nil, err
		}
		items = append(items, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
