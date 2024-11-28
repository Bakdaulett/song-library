package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"song-library/internal/service"
)

type Handler struct {
	SongService *service.SongService
}

func NewHandler(songService *service.SongService) *Handler {
	return &Handler{SongService: songService}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	songs := router.Group("/songs")
	{
		songs.GET("/", h.GetSongs)
		songs.POST("/", h.AddSong)
		songs.GET("/:id", h.GetSongByID)
		songs.PUT("/:id", h.UpdateSong)
		songs.DELETE("/:id", h.DeleteSong)
		songs.GET("/:id/lyrics", h.GetSongLyrics)
		songs.GET("/:id/lyrics/:range", h.GetSongLyricsByRange)
	}

	return router
}
