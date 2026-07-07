package handlers

import(
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type UserHandler struct{
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler{
	return &UserHandler{Repo: repo}
}

func(h *UserHandler) Signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var input struct{
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
		Role string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil{
		writeError(w,http.StatusBadRequest,"invalid request body")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil{
		writeError(w,http.StatusInternalServerError,"could not generate a hash password")
		return
	}
	u := models.User{
		Name: input.Name,
		Email: input.Email,
		PasswordHash: string(hash),
		Role: input.Role,
	}

	err = h.Repo.Insert(&u)
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not insert user")
		return
	}
	writeJSON(w, http.StatusCreated,u)
}