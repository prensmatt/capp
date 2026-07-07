package repository

import(
	"database/sql"
	"errors"

	"ecommerce/internal/models"
)

type UserRepository struct{
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository{
	return &UserRepository{DB: db}
}

func(r *UserRepository) Insert(u *models.User) error{
	query := `INSERT INTO users(name, email, password_hash, role)
	VALUES($1,$2,$3,$4) 
	RETURNING id, created_at`

	err := r.DB.QueryRow(query, u.Name, u.Email, u.PasswordHash, u.Role).Scan(&u.ID,&u.CreatedAt)
	return err
}

func(r *UserRepository) GetByEmail(email string)(*models.User, error){
	var u models.User
	query := `SELECT id,email,name,password_hash,role,created_at
	 FROM users WHERE email=$1`
	err := r.DB.QueryRow(query,email).Scan(&u.ID,&u.Email,&u.Name,&u.PasswordHash,&u.Role,&u.CreatedAt)

	if errors.Is(err,sql.ErrNoRows){
		return nil, models.ErrNotFound
	}
	if err != nil{
		return nil, err
	}
	return &u, nil
}

