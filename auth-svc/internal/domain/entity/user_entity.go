package entity

import "time"

// user entity
type User struct {
	ID          uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Active      bool
	FirstName   string
	LastName    string
	Email       string
	Address     string
	PhoneNumber string
	Password    string
}
