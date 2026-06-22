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


func (r *OrderRepository) GetByID(id int)(*models.Order,error){
	var o models.Order
	query := `SELECT id,user_id,status,total_price,created_at
						FROM orders WHERE id=$1
	`
	err := r.DB.QueryRow(query,id).Scan(&o.ID,&o.UserID,&o.Status,&o.TotalPrice,&o.CreatedAt)
	if errors.Is(err, sql.ErrNoRows){
		return nil,models.ErrNotFound
	}
	if err != nil{
		return nil,err
	}

	itemsQuery := `SELECT id,order_id,product_id,quantity,unit_price
								FROM order_items WHERE order_id=$1
	`
	rows,err := r.DB.Query(itemsQuery,id)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var item models.OrderItem
		if err := rows.Scan(&item.ID,&item.OrderID,
			&item.ProductID,&item.Quantity,&item.UnitPrice,);err != nil{
				return nil, err
			}
			o.Items = append(o.Items,item)
	}
	return &o,rows.Err()
}