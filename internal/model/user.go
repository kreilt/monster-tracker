package model

import "time"

type User struct {
	UserID     int       `json:"user_id"`
	ExternalID string    `json:"-"`
	Nickname   *string   `json:"nickname"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
}
