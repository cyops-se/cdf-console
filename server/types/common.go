package types

import (
	"database/sql"
	"time"
)

type Context struct {
	Cmd     string
	Wdir    string
	Trace   bool
	Version bool
}

type Model struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
}
