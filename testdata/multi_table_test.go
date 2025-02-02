package test

import (
	"time"
)

type users struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsMember  bool
}

type image struct {
	Name      string
	Type      string
	Size      int
	Hidden    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
