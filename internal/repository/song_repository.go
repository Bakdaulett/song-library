package repository

import (
	"context"
	"database/sql"
	"fmt"
	"song-library/internal/models"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) Create(ctx context.Context, song models.Song) (string, error) {
	query := `
		INSERT INTO songs (group_name, song_name, lyrics) 
		VALUES ($1, $2, $3) 
		RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query, song.GroupName, song.SongName, song.Lyrics).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to insert song: %v", err)
	}
	return id, nil
}

func (r *SongRepository) GetByID(ctx context.Context, id string) (*models.Song, error) {
	query := `
		SELECT id, group_name, song_name, lyrics, created_at, updated_at
		FROM songs 
		WHERE id = $1`
	var song models.Song
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.Lyrics,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get song by id: %v", err)
	}
	return &song, nil
}

func (r *SongRepository) GetAll(ctx context.Context, filter string, limit, offset int) ([]models.Song, error) {
	query := `
		SELECT id, group_name, song_name, lyrics, created_at, updated_at
		FROM songs
		WHERE group_name ILIKE $1 OR song_name ILIKE $1
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, "%"+filter+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get songs: %v", err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(
			&song.ID,
			&song.GroupName,
			&song.SongName,
			&song.Lyrics,
			&song.CreatedAt,
			&song.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan song: %v", err)
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (r *SongRepository) Update(ctx context.Context, song models.Song) error {
	query := `
		UPDATE songs
		SET group_name = $1, song_name = $2, lyrics = $3, updated_at = NOW()
		WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, song.GroupName, song.SongName, song.Lyrics, song.ID)
	if err != nil {
		return fmt.Errorf("failed to update song: %v", err)
	}
	return nil
}

func (r *SongRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %v", err)
	}
	return nil
}
