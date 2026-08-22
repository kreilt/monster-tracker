package model

import "time"

type Mark struct {
	MarkID      int        `json:"flavor_user_id"`
	UserID      int        `json:"user_id"`
	FlavorID    int        `json:"flavor_id"`
	StatusTried bool       `json:"status_tried"`
	TriedAt     *time.Time `json:"tried_at"`
	UserPhoto   *string    `json:"user_photo"`
}
