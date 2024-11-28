package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"song-library/internal/service"
	"strconv"
	"strings"
)

// @Summary Get all songs
// @Description Get a list of songs with pagination and optional filters
// @Tags songs
// @Param group query string false "Group filter"
// @Param song query string false "Song filter"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {array} service.Song
// @Failure 400 {object} gin.H{"error": "Invalid page number"}
// @Failure 500 {object} gin.H{"error": "Could not fetch songs"}
// @Router /songs [get]
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

// @Summary Get a song by ID
// @Description Get the details of a song by its ID
// @Tags songs
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} gin.H{"song": service.Song}
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 404 {object} gin.H{"error": "Song not found"}
// @Router /songs/{id} [get]
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

// @Summary Get song lyrics
// @Description Get the lyrics of a song, can be paginated
// @Tags songs
// @Param id path int true "Song ID"
// @Success 200 {object} gin.H{"lyrics": string}
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 500 {object} gin.H{"error": "Could not fetch lyrics"}
// @Router /songs/{id}/lyrics [get]
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

// @Summary Get song lyrics by range
// @Description Get the lyrics of a song by a specific range of verses (e.g., 1-3)
// @Tags songs
// @Param id path int true "Song ID"
// @Param range path string true "Range of verses (e.g., 1-3)"
// @Success 200 {object} gin.H{"lyrics": string}
// @Failure 400 {object} gin.H{"error": "Invalid range format"}
// @Failure 500 {object} gin.H{"error": "Could not fetch lyrics"}
// @Router /songs/{id}/lyrics/{range} [get]
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

// @Summary Add a new song
// @Description Add a new song to the library
// @Tags songs
// @Param song body service.SongRequest true "New song details"
// @Success 201 {object} gin.H{"message": "Song created successfully"}
// @Failure 400 {object} gin.H{"error": "Invalid request body"}
// @Failure 500 {object} gin.H{"error": "Could not add song"}
// @Router /songs [post]
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

// @Summary Update a song
// @Description Update the details of an existing song by ID
// @Tags songs
// @Param id path int true "Song ID"
// @Param song body service.SongRequest true "Updated song details"
// @Success 200 {object} gin.H{"message": "Song updated successfully"}
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 500 {object} gin.H{"error": "Could not update song"}
// @Router /songs/{id} [put]
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

// @Summary Delete a song
// @Description Delete a song by its ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204 {object} gin.H{"message": "Song deleted successfully"}
// @Failure 400 {object} gin.H{"error": "Invalid song ID"}
// @Failure 500 {object} gin.H{"error": "Could not delete song"}
// @Router /songs/{id} [delete]
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
