package handlers

import(
	"errors"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"time"
	"net/http"
	"strconv"
	"encoding/json"

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

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var p models.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil{
		writeError(w, http.StatusBadRequest,"invalid request body")
		return
	}
	err = h.Repo.Insert(&p)
	if err != nil{
		writeError(w, http.StatusInternalServerError, "could not create product")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	limit := 10
	if l := r.URL.Query().Get("limit"); l != ""{
		if parsed,err := strconv.Atoi(l); err == nil{
			limit = parsed
		}
	}
	offset := 0
	if o := r.URL.Query().Get("offset"); o != ""{
		if parsed, err := strconv.Atoi(o); err == nil{
			offset = parsed
		}
	}

	products,err := h.Repo.GetAll(limit,offset)
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not fetch products")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var p models.Product
	id,err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil{
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p.ID = id
	err = h.Repo.Update(&p)
	if errors.Is(err, models.ErrNotFound){
		writeError(w, http.StatusNotFound,"product not found")
		return
	}
	if err != nil{
		writeError(w,http.StatusInternalServerError,"could not update the product")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w, http.StatusBadRequest,"invalid id")
		return
	}
	err = h.Repo.Delete(id)
	if errors.Is(err, models.ErrNotFound){
		writeError(w,http.StatusNotFound,"product not found")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError, "could not delete the product")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) UploadProductImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil{
		writeError(w,http.StatusBadRequest,"invalid id")
		return
	}
	file,header,err := r.FormFile("image")
	if err != nil{
		writeError(w,http.StatusBadRequest,"no image file provided")
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%d%s", id, time.Now().UnixNano(), ext)

	dst, err := os.Create(filepath.Join("static","images",filename))
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not save image")
		return
	}
	defer dst.Close()

	_,err = io.Copy(dst,file)
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not save image")
		return
	}
	imageURL := "/static/images/" + filename
	product, err := h.Repo.GetByID(id)
	if errors.Is(err, models.ErrNotFound){
		writeError(w, http.StatusNotFound,"product not found")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not update product")
		return
	}
	product.ImageURL = imageURL
	err = h.Repo.Update(product)
	if err != nil{
		writeError(w, http.StatusInternalServerError, "could not update product")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"image_url": imageURL})
}