package handlers

import (
	// "encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type SongHandler struct {
	Repo *repository.SongRepository
}

func NewSongHandler(repo *repository.SongRepository) *SongHandler {
	return &SongHandler{Repo: repo}
}

func (h *SongHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	songs, err := h.Repo.GetAll()
	if err != nil {
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not fetch songs")
		return
	}
	writeJSON(w, http.StatusOK, songs)
}

func (h *SongHandler) Upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "could not parse form")
		return
	}

	title := r.FormValue("title")
	artist := r.FormValue("artist")

	if title == "" || artist == "" {
		writeError(w, http.StatusBadRequest, "title and artist are required")
		return
	}

	file, header, err := r.FormFile("song")
	if err != nil {
		writeError(w, http.StatusBadRequest, "no song file provided")
		return
	}
	defer file.Close()

	filename := strconv.FormatInt(int64(header.Size), 10) + "_" + header.Filename
	dst, err := createFile("static/songs/" + filename)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not save song")
		return
	}
	defer dst.Close()

	if _, err := copyFile(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "could not save song")
		return
	}

	s := models.Song{
		Title:   title,
		Artist:  artist,
		FileURL: "/static/songs/" + filename,
	}

	if err := h.Repo.Insert(&s); err != nil {
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not save song")
		return
	}

	writeJSON(w, http.StatusCreated, s)
}

func (h *SongHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			writeError(w, http.StatusNotFound, "song not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not delete song")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}