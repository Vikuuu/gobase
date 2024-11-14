package parser

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
