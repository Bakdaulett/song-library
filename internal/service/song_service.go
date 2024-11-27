package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/repository"
	"time"
)

type SongRequest struct {
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"release_date"`
	Lyrics      string    `json:"lyrics"`
}

type SongService struct {
	SongRepo *repository.SongRepository
}

func NewSongService(songRepo *repository.SongRepository) *SongService {
	return &SongService{SongRepo: songRepo}
}

// GetSongInfoFromAPI fetches song information (e.g., lyrics) from an external API
func (s *SongService) GetSongInfoFromAPI(songName string) (string, error) {
	// Define the API endpoint (example API URL)
	apiURL := fmt.Sprintf("https://api.example.com/songs/%s", songName)

	// Make the HTTP request to the external API
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("error calling external API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from API: %v", err)
	}

	// Unmarshal the response body into a map
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling API response: %v", err)
	}

	// Validate the response and return the lyrics if present
	if description, ok := apiResponse["description"].(string); ok {
		return description, nil
	}

	return "", fmt.Errorf("no description found in API response")
}

// GetSongs retrieves a list of songs with optional filtering and pagination
func (s *SongService) GetSongs(group, song string, page, limit int) ([]models.Song, error) {
	return s.SongRepo.GetSongs(group, song, page, limit)
}

func (s *SongService) GetSongByID(id int) (*models.Song, error) {
	song, err := s.SongRepo.GetSongByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve song with ID %d: %v", id, err)
	}

	return song, nil
}

// AddSong adds a new song to the repository, potentially fetching additional info (e.g., lyrics) from an API
func (s *SongService) AddSong(songRequest SongRequest) error {
	// If no lyrics are provided, try to fetch them from an external API
	if songRequest.Lyrics == "" {
		lyrics, err := s.GetSongInfoFromAPI(songRequest.Song)
		if err != nil {
			return fmt.Errorf("could not fetch lyrics from external API: %v", err)
		}
		songRequest.Lyrics = lyrics
	}

	// Create a new Song model
	song := models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: songRequest.ReleaseDate,
		Lyrics:      songRequest.Lyrics,
	}

	// Save the song in the repository
	return s.SongRepo.AddSong(song)
}

// UpdateSong updates an existing song by its ID
func (s *SongService) UpdateSong(id int, songRequest SongRequest) error {
	// Create the updated song model
	song := models.Song{
		Group:       songRequest.Group,
		Song:        songRequest.Song,
		ReleaseDate: songRequest.ReleaseDate,
		Lyrics:      songRequest.Lyrics,
	}

	// Update the song in the repository
	return s.SongRepo.UpdateSong(id, song)
}

// DeleteSong deletes a song by its ID
func (s *SongService) DeleteSong(id int) error {
	// Delete the song from the repository
	return s.SongRepo.DeleteSong(id)
}
