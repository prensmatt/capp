package handlers

import(
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"
	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type SongHandler struct{
	Repo *repository.SongRepository
}

func NewSongHandler(repo *repository.SongRepository) *SongHandler{
	return &SongHandler{Repo: repo}
}

func(h *SongHandler) GetAll(w http.ResponseWriter, r *http.Request ,ps httprouter.Params){
	songs, err := h.Repo.GetAll()
	if err != nil{
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not fetch songs")
		return
	}
	writeJSON(w, http.StatusOK, songs)
}