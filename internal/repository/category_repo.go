package repository

import(
	"database/sql"
	"errors"

	"ecommerce/internal/models"
)

type CategoryRepository struct{
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB)*CategoryRepository{
	return &CategoryRepository{DB: db}
}

func(r *CategoryRepository) GetAll()([]*models.Category, error){
	rows, err := r.DB.Query(`SELECT id,name,slug FROM categories ORDER BY id ASC`)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next(){
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil{
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, rows.Err()
}

func(r *CategoryRepository) GetByID(id int) (*models.Category ,error){
	var c models.Category
	err:= r.DB.QueryRow(`SELECT id, name, slug FROM categories WHERE id=$1`,id).Scan(&c.ID, &c.Name, &c.Slug)

	if errors.Is(err, sql.ErrNoRows){
		return nil, models.ErrNotFound
	}
	if err != nil{
		return nil, err
	}
		return &c, nil
}
