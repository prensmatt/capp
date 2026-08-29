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

func(r *CategoryRepository) Insert(c *models.Category) error{
	err := r.DB.QueryRow(`INSERT INTO categories(name, slug) VALUES ($1, $2) RETURNING id`,c.Name, c.Slug).Scan(&c.ID)
	return err
}

func(r *CategoryRepository) Update(c *models.Category) error{
	result, err := r.DB.Exec(`UPDATE categories SET name = $1, slug = $2 WHERE id = $3`, c.Name, c.Slug, c.ID)
	if err != nil{
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil{
		return err
	}
	if rows == 0 {
		return models.ErrNotFound
	}
	return nil
}