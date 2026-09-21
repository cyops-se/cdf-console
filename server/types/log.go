package types

import (
	"time"
)

type Log struct {
	Model
	Time        time.Time `json:"time"`
	Source      string    `json:"source"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
