package models

import "time"

type Song struct {
	ID          int        `json:"id"`
	GroupID     int        `json:"-"`
	Group       string     `json:"group"`
	Song        string     `json:"song"`
	ReleaseDate *time.Time `json:"release_date"`
	Lyrics      string     `json:"lyrics"`
	Link        string     `json:"link"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
}
