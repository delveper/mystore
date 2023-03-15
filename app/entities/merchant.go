package entities

import "time"

type Merchant struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Phone     int        `json:"phone"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
