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

