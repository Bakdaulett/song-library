package service

import (
	"fmt"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/repository"
	"strconv"
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

func (s *SongService) ValidateSongRequest(songRequest SongRequest) error {
	if songRequest.Song == "" {
		return fmt.Errorf("song name cannot be empty")
	}
	return nil
}

func (s *SongService) GetSongs(group, song string, page, limit int) ([]models.Song, error) {
	return s.SongRepo.GetSongs(group, song, page, limit)
}

func (s *SongService) GetSongByID(id string) (*models.Song, error) {
	song, err := s.SongRepo.GetSongByID(id)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *SongService) GetSongLyricsWithRange(songID, start, end int) (string, error) {
	song, err := s.SongRepo.GetSongByID(strconv.Itoa(songID))
	if err != nil {
		return "", fmt.Errorf("could not retrieve song with ID %d: %w", songID, err)
	}

	formattedLyrics := strings.ReplaceAll(song.Lyrics, "\\n", "\n")

	verses := strings.Split(strings.TrimSpace(formattedLyrics), "\n\n")
	totalVerses := len(verses)

	if start < 1 || start > totalVerses {
		return "", fmt.Errorf("start index out of range: valid range is 1-%d", totalVerses)
	}

	if end == 0 || end > totalVerses {
		end = totalVerses
	}

	if start > end {
		return "", fmt.Errorf("invalid range: start (%d) cannot be greater than end (%d)", start, end)
	}

	paginatedVerses := verses[start-1 : end]

	return strings.Join(paginatedVerses, "\n\n"), nil
}

func (s *SongService) AddSong(songRequest SongRequest) error {
	if err := s.ValidateSongRequest(songRequest); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	song := models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: &songRequest.ReleaseDate,
		Lyrics:      songRequest.Lyrics,
		Link:        songRequest.Link,
	}

	if err := s.SongRepo.AddSong(song); err != nil {
		return fmt.Errorf("failed to save song: %w", err)
	}
	return nil
}

func (s *SongService) UpdateSong(id int, songRequest SongRequest) error {
	if err := s.ValidateSongRequest(songRequest); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	song := models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: &songRequest.ReleaseDate,
		Lyrics:      songRequest.Lyrics,
		Link:        songRequest.Link,
	}

	if err := s.SongRepo.UpdateSong(id, song); err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}
	return nil
}

func (s *SongService) DeleteSong(id int) error {
	return s.SongRepo.DeleteSong(id)
}
