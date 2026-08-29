package handlers

import(
	"errors"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type CategoryHandler struct{
	Repo *repository.CategoryRepository
}

func NewCategoryHandler(repo *repository.CategoryHandler) *CategoryHandler{
	return &CategoryHandler{Repo: repo}
}

func(h *CategoryHandler) GetAll(w http.ResponseWriter,r *http.Request, ps httprouter.Params){
	categories, err := h.Repo.GetAll()
	if err != nil{
		log.Print(err)
		writeError(w, http.StatusInternalServerError, "could not fetch categories")
		return
	}
	writeJSON(w, http.StatusOK, categories)
}