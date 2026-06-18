package models

import(
	"time"
)

type User struct{
	ID 							int				`json:"id"`
	Email 					string		`json:"email"`
	PasswordHash		string		`json:"-"`
	Name						string		`json:"name"`
	Role						string		`json:"role"`
	CreatedAt				time.Time	`json:"created_at"`
}

type Category struct{
	ID							int				`json:"id"`
	Name						string		`json:"name"`
	Slug						string		`json:"slug"`
}

type Product struct{
	ID							int				`json:"id"`
	Name						string		`json:"name"`
	Stock						int				`json:"stock"`
	CategoryID			int				`json:"category_id"`
	Slug						string		`json:"slug"`
	Description			string		`json:"description"`
	Price						float64		`json:"price"`
	ImageURL				string		`json:"image_url"`
	CreatedAt				time.Time	`json:"created_at"`
}

type Order struct{
	ID							int						`json:"id"`
	UserID					int						`json:"user_id"`
	Status 					string				`json:"status"`
	TotalPrice			float64				`json:"total_price"`
	CreatedAt				time.Time			`json:"created_at"`
	Items						[]OrderItem		`json:"items,omitempty"`
}

type OrderItem struct{
	ID							int						`json:"id"`
	OrderID					int						`json:"order_id"`
	ProductID				int						`json:"product_id"`
	Quantity				int						`json:"quantity"`
	UnitPrice				float64				`json:"unit_price"`
}
