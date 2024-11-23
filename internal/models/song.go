package models

import "time"

type Song struct {
	ID        string    `json:"id"`
	GroupName string    `json:"group"`
	SongName  string    `json:"song"`
	Lyrics    string    `json:"lyrics"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
