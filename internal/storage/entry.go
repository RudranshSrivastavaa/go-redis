package storage

import "time"

type Entry struct {
	Value Value

	CreatedAt time.Time
	UpdatedAt time.Time

	ExpiresAt *time.Time

	LastAccess int64
}