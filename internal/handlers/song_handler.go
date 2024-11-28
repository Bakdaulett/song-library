package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"song-library/internal/service"
	"strconv"
)

// GetSongs handles the GET request to fetch all songs
func (h *Handler) GetSongs(c *gin.Context) {
	group := c.DefaultQuery("group", "")
	song := c.DefaultQuery("song", "")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	songs, err := h.SongService.GetSongs(group, song, page, limit)
	if err != nil {
		log.Printf("Error fetching songs: %v", err)
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
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}
	song, err := h.SongService.GetSongByID(id, page, pageSize)
	if err != nil {
		log.Printf("Error fetching song with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"song": song})

}

// GetSongLyrics handles the GET request to fetch song lyrics with pagination
func (h *Handler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}
	lyrics, err := h.SongService.GetSongLyricsWithPagination(id, page, limit)
	if err != nil {
		if err.Error() == "page out of range" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		} else {
			log.Printf("Error fetching lyrics for song ID %d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch lyrics"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
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
		log.Printf("Error adding song: %v", err)
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
		log.Printf("Error updating song with ID %d: %v", id, err)
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
		log.Printf("Error deleting song with ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete song"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Song deleted successfully"})
}
