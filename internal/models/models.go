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

