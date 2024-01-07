package database

import (
	"context"
	"fmt"
	"os"
	"reflect"

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

func SelectOne[T any](table string, filter map[string]string, result interface{}) error {
	query := psql.Select("*").From(table)

	for key, value := range filter {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	row := db.QueryRow(context.Background(), sql, args...)

	// Get type information of the result
	valueType := reflect.TypeOf(result).Elem()

	// Create a slice to hold the pointers to fields
	pointers := make([]interface{}, valueType.NumField())

	// Iterate over the fields of the struct and create pointers
	for i := 0; i < valueType.NumField(); i++ {
		pointers[i] = reflect.New(valueType.Field(i).Type).Interface()
	}

	// Pass the pointers as variadic arguments to Scan
	if err := row.Scan(pointers...); err != nil {
		return err
	}

	// Assign the scanned values from pointers to the result struct
	resultValue := reflect.ValueOf(result).Elem()
	for i := 0; i < valueType.NumField(); i++ {
		resultValue.Field(i).Set(reflect.ValueOf(pointers[i]).Elem())
	}

	return nil
}

func SelectMany[T any](table string, filter map[string]string, result *[]T) error {
	query := squirrel.Select("*").From(table)

	for key, value := range filter {
		query = query.Where(squirrel.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Create a slice to hold the results
	var users []T

	// Get type information for the slice elements
	elemType := reflect.TypeOf(users).Elem()

	// Iterate through the rows and create a new instance of T for each row
	for rows.Next() {
		// Create a new instance of T
		newItem := reflect.New(elemType).Interface()

		// Scan the row into the newly created instance
		if err := rows.Scan(newItem); err != nil {
			return err
		}

		// Append the scanned item to the slice
		users = append(users, *newItem.(*T))
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return err
	}

	// Assign the populated slice to the result pointer
	*result = users

	return nil
}
