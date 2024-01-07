package models

import (
	"database/sql"
	"encoding/json"
)

// We have to use sql.NullString cause SQL and Go treat NULL differently
// which prevents us from mapping nil values to a struct

// Local wrapper of sql.NullString that implements
// required interfaces for struct mapping from SQL
type NullString struct {
	sql.NullString
}

// Implement the MarshalJSON interface on our wrapper type
func (s NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String)
}
