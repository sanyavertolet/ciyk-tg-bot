package database

import "time"

type Checkpoint struct {
	ID        uint
	Line      int
	CreatedAt time.Time
}
