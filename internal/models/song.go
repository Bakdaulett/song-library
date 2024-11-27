package models

import "time"

type Song struct {
	ID          int        `json:"id"`
	GroupID     int        `json:"group_id"`
	GroupName   string     `json:"group_name"`
	Song        string     `json:"song"`
	ReleaseDate *time.Time `json:"release_date"`
	Lyrics      string     `json:"lyrics"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
