package models

import (
	"database/sql"
	"encoding/json"

	"gorm.io/gorm"
)

var OpenDB *gorm.DB

type NullString struct {
	sql.NullString
}

// parse null string on model
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`null`), nil
}

func NullStringInput(s string) NullString {
	if len(s) == 0 {
		return NullString{sql.NullString{`null`, false}}
	}
	return NullString{sql.NullString{s, true}}
}
