package main

import(
	"net/http"
	"log"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"ecommerce/internal/handlers"
	"ecommerce/internal/repository"
	"ecommerce/internal/config"
)

func main(){

	if err := godotenv.Load(); err != nil{
		log.Println("no .env file found, relying on real environment variables")
	}

	cfg, err := config.Load()
	if err != nil{
		log.Fatal(err)
	}

	db, err := repository.NewDB(cfg.DBURL)
	if err != nil{
		log.Fatal(err)
	}
	
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	productHandler := handlers.NewProductHandler(productRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo)

	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/products", productHandler.GetAllProducts)
	router.GET("/products/:id", productHandler.GetProduct)
	router.POST("/products", productHandler.CreateProduct)
	router.PUT("/products/:id", productHandler.UpdateProduct)
	router.DELETE("/products/:id", productHandler.DeleteProduct)
	router.POST("/products/:id/image", productHandler.UploadProductImage)

	router.GET("/orders", orderHandler.GetAllOrders)
	router.GET("/orders/:id", orderHandler.GetOrder)
	router.POST("/orders", orderHandler.CreateOrder)
	router.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)

	log.Printf("Server starting on: %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}