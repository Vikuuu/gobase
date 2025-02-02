package testdata

import "time"

type User struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsMember  bool
}

type Image struct {
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsMember  bool
}

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	InStock     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Order struct {
	OrderID    string
	UserID     int
	ProductID  string
	Quantity   int
	TotalPrice float64
	OrderDate  time.Time
	Delivered  bool
}

type Review struct {
	ReviewID   string
	ProductID  string
	UserID     int
	Rating     int
	Comment    string
	ReviewDate time.Time
}

type Category struct {
	CategoryID string
	Name       string
	ParentID   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
