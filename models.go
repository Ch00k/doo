package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type ModelBase struct {
	ID        uint  `gorm:"primarykey"`
	CreatedAt int64 `gorm:"autoCreateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

type Entry struct {
	ModelBase
	CompletedAt CompletedAt
	// TODO: This means that on PUT this will be required too.
	// It is ok now that this is the only field that cane be updated,
	// but this needs a custom validator when more fields are added
	Text     string `binding:"required"`
	Comments []Comment
}

type Comment struct {
	ModelBase
	EntryID uint
	Text    string `binding:"required"`
}

type CompletedAt sql.NullInt64

// Scan implements the Scanner interface.
func (n *CompletedAt) Scan(value interface{}) error {
	return (*sql.NullInt64)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n CompletedAt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

func (n CompletedAt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

func (n *CompletedAt) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int64)
	if err == nil {
		n.Valid = true
	}
	return err
}
