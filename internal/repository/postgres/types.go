package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONB is a generic wrapper around any Go type that implements driver.Valuer
// and sql.Scanner for seamless reading and writing of PostgreSQL JSONB columns.
type JSONB[T any] struct {
	Data T
}

func (j JSONB[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j *JSONB[T]) Scan(src any) error {
	if src == nil {
		return nil
	}
	var bytes []byte
	switch v := src.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONB: %T", src)
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, &j.Data)
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func toNullStringPointer(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func fromNullStringPointer(ns sql.NullString) *string {
	if ns.Valid && ns.String != "" {
		val := ns.String
		return &val
	}
	return nil
}
