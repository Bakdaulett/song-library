package handlers

import (
	"encoding/json"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/repository"
	"strconv"

	"github.com/gorilla/mux"
)

type SongHandler struct {
	repo *repository.SongRepository
}

func NewSongHandler(repo *repository.SongRepository) *SongHandler {
	return &SongHandler{repo: repo}
}

func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(r.Context(), song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *SongHandler) GetSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	song, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if song == nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(song)
}

func (h *SongHandler) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	songs, err := h.repo.GetAll(r.Context(), filter, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	song.ID = mux.Vars(r)["id"]

	if err := h.repo.Update(r.Context(), song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
