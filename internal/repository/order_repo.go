package repository

import(
	"database/sql"
	"errors"

	"ecommerce/internal/models"
)

type OrderRepository struct{
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository{
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(o *models.Order) error{
	tx,err := r.DB.Begin()
	if err != nil{
		return err
	}
	defer tx.Rollback()

	query :=`INSERT INTO orders (user_id,status,total_price)
					VALUES($1,$2,$3)
					RETURNING id, created_at
	`
	err = tx.QueryRow(query,o.UserID,o.Status,o.TotalPrice).Scan(&o.ID,&o.CreatedAt)
	if err != nil{
		return err
	}

	for _, item := range o.Items{
		query = `INSERT INTO order_items (order_id,product_id,quantity,unit_price) VALUES($1,$2,$3,$4)`
		_,err = tx.Exec(query,o.ID,item.ProductID,item.Quantity,item.UnitPrice)
		if err != nil{
			return err
		}

		query = `UPDATE products SET stock=stock-$1 WHERE id=$2 AND stock>=$1`
		result, err := tx.Exec(query,item.Quantity, item.ProductID)

		if err != nil{
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil{
			return err
		}
		if rows==0{
			return models.ErrInsufficientStock
		}
	}
	return tx.Commit()
}