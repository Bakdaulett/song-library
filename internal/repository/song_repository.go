package repository

import (
	"database/sql"
	"fmt"
	"log"
	"song-library/internal/models"
)

type SongRepository struct {
	DB *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{DB: db}
}

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

func (r *SongRepository) getOrCreateGroupID(groupName string) (int, error) {
	var groupID int
	query := `SELECT id FROM groups WHERE name = $1`
	err := r.DB.QueryRow(query, groupName).Scan(&groupID)
	if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
		err = r.DB.QueryRow(insertQuery, groupName).Scan(&groupID)
		if err != nil {
			return 0, fmt.Errorf("failed to create group: %w", err)
		}
	} else if err != nil {
		return 0, fmt.Errorf("failed to fetch group ID: %w", err)
	}
	return groupID, nil
}

func (r *SongRepository) AddSong(song models.Song) error {
	groupID, err := r.getOrCreateGroupID(song.Group)
	if err != nil {
		return fmt.Errorf("failed to get or create group ID: %w", err)
	}

	query := `
		INSERT INTO songs (group_id, song, release_date, lyrics, link) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err = r.DB.Exec(query, groupID, song.Song, song.ReleaseDate, song.Lyrics, song.Link)
	if err != nil {
		return fmt.Errorf("failed to insert song: %w", err)
	}

	log.Printf("Successfully added song %q by group %q", song.Song, song.Group)
	return nil
}

func (r *SongRepository) GetSongByID(songID string) (*models.Song, error) {
	query := `
        SELECT s.id, g.name, s.song, s.release_date, s.lyrics, s.link
        FROM songs s
        JOIN groups g ON s.group_id = g.id
        WHERE s.id = $1
    `
	var song models.Song
	err := r.DB.QueryRow(query, songID).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Lyrics, &song.Link)
	if err != nil {
		return nil, err
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("song not found")
		}
		return nil, fmt.Errorf("error fetching song: %w", err)
	}
	return &song, nil
}

func (r *SongRepository) UpdateSong(id int, song models.Song) error {
	groupID, err := r.getOrCreateGroupID(song.Group)
	if err != nil {
		return fmt.Errorf("failed to get or create group ID: %w", err)
	}

	query := `
        UPDATE songs
        SET group_id = $1, song = $2, release_date = $3, lyrics = $4, link = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6
    `
	_, err = r.DB.Exec(query, groupID, song.Song, song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}

	log.Printf("Successfully updated song with ID %d", id)
	return nil
}

func (r *SongRepository) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	log.Printf("Successfully deleted song with ID %d", id)
	return nil
}
