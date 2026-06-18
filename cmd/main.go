package main

import(
	"fmt"
	"os"
	"net/http"
	"log"
)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port="8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health",func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w, "OK")
	})

	log.Printf("Server starting on: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}