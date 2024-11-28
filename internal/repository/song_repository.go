package repository

import (
	"database/sql"
	"fmt"
	"log"
	"song-library/internal/models"
	"strings"
)

type SongRepository struct {
	DB *sql.DB
}

// NewSongRepository creates a new SongRepository instance
func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{DB: db}
}

// GetSongs retrieves songs with pagination, based on group and song name filters
func (r *SongRepository) GetSongs(group, song string, page, limit int) ([]models.Song, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("page and limit must be greater than 0")
	}
	offset := (page - 1) * limit
	query := `
        SELECT s.id, g.name, s.song, s.release_date, s.lyrics, s.link
        FROM songs s
        JOIN groups g ON s.group_id = g.id
        WHERE g.name ILIKE $1 AND s.song ILIKE $2
        ORDER BY s.id
        LIMIT $3 OFFSET $4
    `
	rows, err := r.DB.Query(query, "%"+group+"%", "%"+song+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()
	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Lyrics, &song.Link); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		songs = append(songs, song)
	}
	return songs, nil
}

// AddSong adds a new song to the database and handles adding a group if necessary
func (r *SongRepository) AddSong(song models.Song) error {
	groupID, err := r.getOrCreateGroupID(song.GroupID)
	if err != nil {
		return fmt.Errorf("failed to get or create group ID: %w", err)
	}
	query := `INSERT INTO songs (group_id, song, release_date, lyrics, link) VALUES ($1, $2, $3, $4, $5)`
	_, err = r.DB.Exec(query, groupID, song.Song, song.ReleaseDate, song.Lyrics, song.Link)
	if err != nil {
		return fmt.Errorf("failed to insert song: %w", err)
	}
	log.Printf("Successfully added song %q by group %q", song.Song, song.GroupID)
	return nil
}

func (r *SongRepository) GetSongByID(id int, page int, pageSize int) (*models.Song, error) {
	query := `
        SELECT s.id, g.name, s.song, s.release_date, s.lyrics, s.link
        FROM songs s
        JOIN groups g ON s.group_id = g.id
        WHERE s.id = $1
    `
	var song models.Song
	err := r.DB.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Lyrics, &song.Link)
	if err != nil {
		return nil, err
	}
	couplets := strings.Split(song.Lyrics, "\n\n")
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(couplets) {
		endIndex = len(couplets)
	}
	song.Lyrics = strings.Join(couplets[startIndex:endIndex], "\n\n")
	return &song, nil
}

// UpdateSong updates an existing song in the database by its ID
func (r *SongRepository) UpdateSong(id int, song models.Song) error {
	query := `
        UPDATE songs
        SET song = $1, release_date = $2, lyrics = $3, link = $4, updated_at = CURRENT_TIMESTAMP
        WHERE id = $5
    `
	_, err := r.DB.Exec(query, song.Song, song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}
	log.Printf("Successfully updated song with ID %d", id)
	return nil
}

// DeleteSong deletes a song by its ID from the database
func (r *SongRepository) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	log.Printf("Successfully deleted song with ID %d", id)
	return nil
}

// getOrCreateGroupID retrieves the ID of a group, creating it if it doesn't exist
func (r *SongRepository) getOrCreateGroupID(groupName int) (int, error) {
	var groupID int
	err := r.DB.QueryRow(`INSERT INTO groups (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id`, groupName).Scan(&groupID)
	if err == sql.ErrNoRows {
		err = r.DB.QueryRow(`SELECT id FROM groups WHERE name = $1`, groupName).Scan(&groupID)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to retrieve or create group ID for %q: %w", groupName, err)
	}

	return groupID, nil
}
