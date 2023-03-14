package entities

import "time"

type Customer struct {
	ID        int        `json:"id"`
	FullName  string     `json:"full_name"`
	Phone     int        `json:"phone"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
