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
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not fetch categories")
		return
	}
	writeJSON(w, http.StatusOK, categories)
}

func(h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	category, err := h.Repo.GetByID(id)
	if errors.Is(err, models.ErrNotFound){
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	if err != nil{
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not fetch category")
		return
	}
	writeJSON(w, http.StatusOK, category)
}

func(h *http.CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request, ps http.Params){
	var input struct{
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil{
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(input.Name) == 0 {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if len(input.Slug) == 0 {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	c := models.Category{
		Name: input.Name,
		Slug: input.Slug,
	}

	err = h.repo.Insert(&c)
	if err != nil{
		writeError(w, http.StatusInternalServerError, "could not create category")
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c := models.Category{
		ID:   id,
		Name: input.Name,
		Slug: input.Slug,
	}

	err = h.Repo.Update(&c)
	if errors.Is(err, models.ErrNotFound) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	
	if err != nil {
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not update category")
		return
	}
	writeJSON(w, http.StatusOK, c)
}