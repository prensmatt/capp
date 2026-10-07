package main

import(
	"fmt"
	"os"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil{
		log.Println("no .env file found")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == ""{
		log.Fatal("JWT_SECRET is not set")
	}
	claims := jwt.MapClaims{
		"user_id": 1,
		"role": "admin",
		"exp": time.Now().Add(24*time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(signed)
}