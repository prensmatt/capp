package handlers

import(
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"ecommerce/internal/repository"
	"ecommerce/internal/models"
)

type ProductHandler struct{
	Repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler{
	return &ProductHandler{Repo: repo}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id,err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	product, err := h.Repo.GetByID(id)
	if errors.Is(err, models.ErrNotFound){
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError,"something went wrong")
		return
	}
	writeJSON(w, http.StatusOK, product)
}