package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"song-library/internal/service"
	"strconv"
	"strings"
)

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
	song, err := h.SongService.GetSongByID(id)
	if err != nil {
		log.Printf("Error fetching song with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"song": song})

}

func (h *Handler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	lyrics, err := h.SongService.GetSongLyricsWithRange(id, 1, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch lyrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}

func (h *Handler) GetSongLyricsByRange(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	rangeStr := c.Param("range")

	var start, end int
	if strings.Contains(rangeStr, "-") {
		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range format"})
			return
		}

		start, err = strconv.Atoi(parts[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start range"})
			return
		}

		end, err = strconv.Atoi(parts[1])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end range"})
			return
		}
	} else {
		start, err = strconv.Atoi(rangeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range value"})
			return
		}
		end = start
	}

	lyrics, err := h.SongService.GetSongLyricsWithRange(id, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch lyrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lyrics": lyrics})
}

func (h *Handler) AddSong(c *gin.Context) {
	var songRequest service.SongRequest

	if err := c.ShouldBindJSON(&songRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.SongService.AddSong(songRequest); err != nil {
		log.Printf("Error adding song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully"})
}

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

	if err := h.SongService.UpdateSong(id, songRequest); err != nil {
		log.Printf("Error updating song with ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

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
