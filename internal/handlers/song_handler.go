package handlers

import (
	"encoding/json"
	"net/http"
	"song-library/internal/service"
	"song-library/pkg/logger"
	"strconv"

	"github.com/gorilla/mux"
)

type SongHandler struct {
	SongService *service.SongService
}

func NewSongHandler(songService *service.SongService) *SongHandler {
	return &SongHandler{
		SongService: songService,
	}
}

func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	songs, err := h.SongService.GetSongs(group, song, page, limit)
	if err != nil {
		logger.Error("Error fetching songs: " + err.Error())
		http.Error(w, "Error fetching songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		logger.Error("Error encoding response: " + err.Error())
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song service.SongRequest

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		logger.Error("Error decoding request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.SongService.AddSong(song); err != nil {
		logger.Error("Error adding song: " + err.Error())
		http.Error(w, "Error adding song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songIDStr := vars["id"]

	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var song service.SongRequest

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		logger.Error("Error decoding request body: " + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.SongService.UpdateSong(songID, song); err != nil {
		logger.Error("Error updating song: " + err.Error())
		http.Error(w, "Error updating song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songIDStr := vars["id"]

	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	if err := h.SongService.DeleteSong(songID); err != nil {
		logger.Error("Error deleting song: " + err.Error())
		http.Error(w, "Error deleting song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	songHandler := NewSongHandler(service.NewSongService())

	router.HandleFunc("/songs", songHandler.GetSongs).Methods(http.MethodGet)
	router.HandleFunc("/songs", songHandler.AddSong).Methods(http.MethodPost)
	router.HandleFunc("/songs/:id", songHandler.UpdateSong).Methods(http.MethodPut)
	router.HandleFunc("/songs/:id", songHandler.DeleteSong).Methods(http.MethodDelete)

	return router
}
