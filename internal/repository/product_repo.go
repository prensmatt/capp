package repository

import(
	"database/sql"
	"errors"

	"ecommerce/internal/models"
)

type ProductRepository struct{
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository{
	return ProductRepository{DB: db}
}

func (r *ProductRepository) Insert(p *models.Product) error{
	query := `INSERT INTO products(name,slug,description,price,stock,category_id,image_url)
						VALUES($1,$2,$3,$4,$5,$6,$7)
						RETURNING(id,created_at)
	`

	err := r.DB.QueryRow(query,p.Name,p.Slug,p.Description,p.Price,p.Stock,p.CategoryID,p.ImageURL).Scan(&p.ID,&p.CreatedAt)
	return err
}

func (r *ProductRepository) GetByID(id int) (*models.Product,error){
	query := `SELECT id,name,slug,description,price,stock,category_id,image_url,created_at
						FROM products WHERE id=$1
	`
	err := r.DB.QueryRow(query, id).Scan(&p.ID,&p.Name,&p.Slug,&p.Description,&p.Price,&p.Stock,&p.CategoryID,&p.ImageURL,&p.CreatedAt)
	if errors.Is(err,sql.ErrNoRows){
		return nil, models.ErrNotFound
	}
	if err != nil{
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) GetAll(limit, offset int)([]*models.Product,error){
	query := `SELECT id,name,slug,description,price,stock,category_id,image_url,created_at
						FROM products
						ORDER BY created_at DESC
						LIMIT $1 OFFSET $2
	`
	rows, err := r.DB.Query(query,limit,offset)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next(){
		var p models.Product
		if err := rows.Scan(&p.ID,&p.Name,&p.Slug,&p.Description,
			&p.Price,&p.Stock,&p.CategoryID,&p.ImageURL,&p.CreatedAt,);err != nil{
				return nil, err
			}
			products = append(products,&p)
	}
	return products,rows.Err()
}

func (r *ProductRepository) Update(p *models.Product) error{

}

func (r *ProductRepository) Delete(id int) error{

}