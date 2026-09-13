package handlers

import(
	"net/http"
	"errors"
	"encoding/json"
	"strconv"
	"log"

	"github.com/julienschmidt/httprouter"
	"ecommerce/internal/repository"
	"ecommerce/internal/models"
)

type OrderHandler struct{
	Repo *repository.OrderRepository
	Secret string
}

func NewOrderHandler(repo *repository.OrderRepository, secret string) *OrderHandler{
	return &OrderHandler{Repo: repo, Secret: secret}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Items []models.OrderItem `json:"items"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(input.Items) == 0 {
		writeError(w, http.StatusBadRequest, "order must have at least one item")
		return
	}

	for _, item := range input.Items {
		if item.ProductID <= 0 {
			writeError(w, http.StatusBadRequest, "valid product_id is required for each item")
			return
		}
		if item.Quantity <= 0 {
			writeError(w, http.StatusBadRequest, "quantity must be greater than 0")
			return
		}
		if item.UnitPrice <= 0 {
			writeError(w, http.StatusBadRequest, "unit_price must be greater than 0")
			return
		}
	}

	// get user_id from JWT 
	userID, err := getUserIDFromToken(r, h.Secret)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or missing token")
		return
	}

	o := models.Order{
		UserID: userID,
		Status: "pending",
		Items:  input.Items,
	}

	err = h.Repo.Create(&o)
	if errors.Is(err, models.ErrInsufficientStock) {
		writeError(w, http.StatusUnprocessableEntity, "insufficient stock")
		return
	}
	if err != nil {
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not create order")
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

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	limit := 10
	if l := r.URL.Query().Get("limit"); l != ""{
		if parsed, err := strconv.Atoi(l); err == nil{
			limit = parsed
		}
	}
	offset := 0
	if o := r.URL.Query().Get("offset"); o != ""{
		if parsed, err := strconv.Atoi(o); err == nil{
			offset = parsed
		}
	}
	orders, err := h.Repo.GetAll(limit,offset)
	if err != nil{
		writeError(w, http.StatusInternalServerError, "could not fetch orders")
		return
	}
	writeJSON(w,http.StatusOK,orders)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w,http.StatusBadRequest,"invalid id")
		return
	}
	var body struct{
		Status string `json:"status"`
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil{
		writeError(w,http.StatusBadRequest,"invalid request body")
		return
	}
	err = h.Repo.UpdateStatus(id, body.Status)
	if errors.Is(err, models.ErrNotFound){
		writeError(w,http.StatusNotFound,"order not found")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not update status")
		return
	}
	writeJSON(w, http.StatusOK,map[string]string{"status":body.Status})
}


func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := getUserIDFromToken(r, h.Secret)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or missing token")
		return
	}

	orders, err := h.Repo.GetByUserID(userID)
	if err != nil {
		log.Println(err)
		writeError(w, http.StatusInternalServerError, "could not fetch orders")
		return
	}

	writeJSON(w, http.StatusOK, orders)
}