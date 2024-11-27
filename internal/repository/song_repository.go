package repository

import (
	"database/sql"
	"fmt"
	"song-library/internal/models"
)

type SongRepository struct {
	DB *sql.DB
}

// NewSongRepository creates a new SongRepository instance
func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{DB: db}
}

// GetSongs retrieves songs with pagination, based on the group and song name filters
func (r *SongRepository) GetSongs(group, song string, page, limit int) ([]models.Song, error) {
	query := `
		SELECT s.id, g.name, s.song, s.release_date, s.lyrics
		FROM songs s
		JOIN groups g ON s.group_id = g.id
		WHERE g.name ILIKE $1 AND s.song ILIKE $2
		LIMIT $3 OFFSET $4
	`
	rows, err := r.DB.Query(query, "%"+group+"%", "%"+song+"%", limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
		if err := rows.Scan(&s.ID, &s.Group, &s.Song, &s.ReleaseDate, &s.Lyrics); err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	return songs, nil
}

// AddSong adds a new song to the database and handles adding a group if necessary
func (r *SongRepository) AddSong(song models.Song) error {
	// First, try to add the group if it doesn't exist
	var groupID int
	err := r.DB.QueryRow(`INSERT INTO groups (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`, song.Group).
		Scan(&groupID)
	if err != nil {
		return fmt.Errorf("failed to add group: %v", err)
	}

	if groupID == 0 {
		// If the group already exists, get its ID
		err = r.DB.QueryRow(`SELECT id FROM groups WHERE name = $1`, song.Group).Scan(&groupID)
		if err != nil {
			return fmt.Errorf("failed to fetch group ID: %v", err)
		}
	}

	// Add the song with the obtained group ID
	query := `INSERT INTO songs (group_id, song, release_date, lyrics) VALUES ($1, $2, $3, $4)`
	_, err = r.DB.Exec(query, groupID, song.Song, song.ReleaseDate, song.Lyrics)
	return err
}

// GetSongByID retrieves a song by its ID
func (r *SongRepository) GetSongByID(id int) (*models.Song, error) {
	var song models.Song

	query := `
		SELECT s.id, g.name, s.song, s.release_date, s.lyrics
		FROM songs s
		JOIN groups g ON s.group_id = g.id
		WHERE s.id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Lyrics)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("song with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching song by ID: %v", err)
	}

	return &song, nil
}

// UpdateSong updates an existing song in the database by its ID
func (r *SongRepository) UpdateSong(id int, song models.Song) error {
	query := `
		UPDATE songs
		SET song = $1, release_date = $2, lyrics = $3
		WHERE id = $4
	`
	_, err := r.DB.Exec(query, song.Song, song.ReleaseDate, song.Lyrics, id)
	return err
}

// DeleteSong deletes a song by its ID from the database
func (r *SongRepository) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
