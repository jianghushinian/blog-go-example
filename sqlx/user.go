package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID        int
	Name      sql.NullString `json:"username"`
	Email     string
	Age       int
	Birthday  time.Time
	Salary    Salary
	CreatedAt time.Time `db:"created_at"` // db tag default field to lower: createdat
	UpdatedAt time.Time `db:"updated_at"` // db tag default field to lower: updatedat
}

type Salary struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// Scan implements sql.Scanner, use custom types in *sql.Rows.Scan
func (s *Salary) Scan(src any) error {
	if src == nil {
		return nil
	}

	var buf []byte
	switch v := src.(type) {
	case []byte:
		buf = v
	case string:
		buf = []byte(v)
	default:
		return fmt.Errorf("invalid type: %T", src)
	}

	err := json.Unmarshal(buf, s)
	return err
}

// Value implements driver.Valuer, use custom types in Query/QueryRow/Exec
func (s Salary) Value() (driver.Value, error) {
	v, err := json.Marshal(s)
	return string(v), err
}
