package handlers

import(
	"encoding/json"
	"net/http"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type UserHandler struct{
	Repo *repository.UserRepository
	Secret string
}

func NewUserHandler(repo *repository.UserRepository, secret string) *UserHandler{
	return &UserHandler{Repo: repo, Secret: secret}
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

func(h *UserHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var input struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil{
		writeError(w,http.StatusBadRequest,"invalid request body")
		return
	}
	user, err := h.Repo.GetByEmail(input.Email)
	if errors.Is(err, models.ErrNotFound){
		writeError(w,http.StatusUnauthorized,"invalid credentials")
		return
	}
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not find user")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(input.Password))
	if err != nil{
		writeError(w,http.StatusUnauthorized,"invalid credentials")
		return
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(24*time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.Secret))
	if err != nil{
		writeError(w, http.StatusInternalServerError,"could not sign password")
		return
	}
	writeJSON(w,http.StatusOK,map[string]string{"token":signed})
}