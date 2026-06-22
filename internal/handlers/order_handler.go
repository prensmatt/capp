package handlers

import(
	"net/http"
	"errors"
	"encoding/json"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"ecommerce/internal/repository"
	"ecommerce/internal/models"
)

type OrderHandler struct{
	Repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler{
	return &OrderHandler{Repo: repo}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var o models.Order
	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil{
		writeError(w, http.StatusBadRequest,"invalid request body")
		return
	}
	err = h.Repo.Create(&o)
	if errors.Is(err, models.ErrInsufficientStock){
		writeError(w, http.StatusUnprocessableEntity, "insufficient stock")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not create order")
		return
	}
	writeJSON(w, http.StatusCreated, o)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w,http.StatusBadRequest,"invalid id")
		return
	}
	order, err := h.Repo.GetByID(id)
	if errors.Is(err, models.ErrNotFound){
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not fetch the order")
		return
	}
	writeJSON(w, http.StatusOK, order)
}