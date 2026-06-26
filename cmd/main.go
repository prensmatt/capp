package main

import(
	"net/http"
	"log"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"ecommerce/internal/handlers"
	"ecommerce/internal/repository"
	"ecommerce/internal/config"
	"ecommerce/internal/middleware"
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

	protect := middleware.Auth(cfg.JWTSecret)

	wrap := func(h httprouter.Handle) httprouter.Handle{
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
			protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
				h(w,r,ps)
			})).ServeHTTP(w,r)
		}
	}

	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/products", productHandler.GetAllProducts)
	router.GET("/products/:id", productHandler.GetProduct)

	router.POST("/products", wrap(productHandler.CreateProduct))
	router.PUT("/products/:id", wrap(productHandler.UpdateProduct))
	router.DELETE("/products/:id", wrap(productHandler.DeleteProduct))
	router.POST("/products/:id/image", wrap(productHandler.UploadProductImage))

	router.GET("/orders", wrap(orderHandler.GetAllOrders))
	router.GET("/orders/:id", wrap(orderHandler.GetOrder))
	router.PATCH("/orders/:id/status", wrap(orderHandler.UpdateOrderStatus))

	router.POST("/orders", orderHandler.CreateOrder)

	log.Printf("Server starting on: %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}