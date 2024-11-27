package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"song-library/internal/service"
	"strconv"
)

// GetSongs handles the GET request to fetch all songs
func (h *Handler) GetSongs(c *gin.Context) {
	group := c.DefaultQuery("group", "")
	song := c.DefaultQuery("song", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	songs, err := h.SongService.GetSongs(group, song, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// GetSongByID handles the GET request to fetch a song by its ID
func (h *Handler) GetSongByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song, err := h.SongService.GetSongByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// AddSong handles the POST request to add a new song
func (h *Handler) AddSong(c *gin.Context) {
	var songRequest service.SongRequest
	if err := c.ShouldBindJSON(&songRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the release date
	if songRequest.ReleaseDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date"})
		return
	}

	err := h.SongService.AddSong(songRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully"})
}

// UpdateSong handles the PUT request to update a song by its ID
func (h *Handler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var songRequest service.SongRequest
	if err := c.ShouldBindJSON(&songRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the release date
	if songRequest.ReleaseDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date"})
		return
	}

	err = h.SongService.UpdateSong(id, songRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// DeleteSong handles the DELETE request to remove a song by its ID
func (h *Handler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	err = h.SongService.DeleteSong(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete song"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Song deleted successfully"})
}
