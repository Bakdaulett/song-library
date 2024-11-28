package service

import (
	"fmt"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/repository"
	"strings"
	"time"
)

type SongRequest struct {
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"release_date"`
	Lyrics      string    `json:"lyrics"`
	Link        string    `json:"link"`
}

type SongService struct {
	SongRepo   *repository.SongRepository
	HTTPClient *http.Client
}

func NewSongService(songRepo *repository.SongRepository) *SongService {
	return &SongService{
		SongRepo:   songRepo,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ValidateSongRequest validates the fields of a SongRequest
func (s *SongService) ValidateSongRequest(songRequest SongRequest) error {
	if songRequest.Song == "" {
		return fmt.Errorf("song name cannot be empty")
	}
	return nil
}

// GetSongs retrieves a list of songs with optional filtering and pagination
func (s *SongService) GetSongs(group, song string, page, limit int) ([]models.Song, error) {
	return s.SongRepo.GetSongs(group, song, page, limit)
}

// GetSongByID retrieves a song by its ID
func (s *SongService) GetSongByID(id int, page int, pageSize int) (*models.Song, error) {
	song, err := s.SongRepo.GetSongByID(id, page, pageSize)
	if err != nil {
		return nil, err
	}

	return song, nil
}

// GetSongLyricsWithPagination fetches song lyrics with pagination by verses
func (s *SongService) GetSongLyricsWithPagination(songID, page, limit int) (string, error) {
	song, err := s.SongRepo.GetSongByID(songID, page, limit)
	if err != nil {
		return "", fmt.Errorf("could not retrieve song with ID %d: %w", songID, err)
	}
	verses := strings.Split(song.Lyrics, "\n\n")
	totalVerses := len(verses)
	start := (page - 1) * limit
	end := start + limit
	if start >= totalVerses {
		return "", fmt.Errorf("page out of range")
	}
	if end > totalVerses {
		end = totalVerses
	}
	paginatedVerses := verses[start:end]
	return strings.Join(paginatedVerses, "\n\n"), nil
}

// AddSong adds a new song to the repository, potentially fetching additional info (e.g., lyrics) from an API
func (s *SongService) AddSong(songRequest SongRequest) error {
	// Validate the song request
	if err := s.ValidateSongRequest(songRequest); err != nil {
		return err
	}

	// Create a new Song model
	song := models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: &songRequest.ReleaseDate,
		Lyrics:      songRequest.Lyrics,
	}

	// Save the song in the repository
	return s.SongRepo.AddSong(song)
}

// UpdateSong updates an existing song by its ID
func (s *SongService) UpdateSong(id int, songRequest SongRequest) error {
	// Validate the song request
	if err := s.ValidateSongRequest(songRequest); err != nil {
		return err
	}

	// Create the updated song model
	song := models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: &songRequest.ReleaseDate,
		Lyrics:      songRequest.Lyrics,
		Link:        songRequest.Link,
	}

	// Update the song in the repository
	return s.SongRepo.UpdateSong(id, song)
}

// DeleteSong deletes a song by its ID
func (s *SongService) DeleteSong(id int) error {
	// Delete the song from the repository
	return s.SongRepo.DeleteSong(id)
}
