/*
* This file defines the input and output
* expected in the gobase v0.0.1
 */
package gobase

import (
	"time"
)

// Input struct
type users struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	IsMember  bool
	NewField  string
}

// Output SQL file

/*
* CREATE TABLE users (
* 	id INT,
* 	created_at TIMESTAMP,
* 	updated_at TIMESTAMP,
* 	name TEXT
* );
 */
