package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/gorilla/handlers"
)

// our main function
func main() {
	fmt.Println("Starting...")

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	compose := os.Getenv("indocker")
	if compose == "dockercompose" {
		fmt.Println("Sleeping for 5 to wait for DB")
		time.Sleep(time.Second * 5)
	}

	ReadyDB()

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.Use(JwtAuthentication)

	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/quote", QuoteResponse).Methods("GET")

	p := fmt.Sprintf(":%v", port)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		Debug:            true,
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	handler := c.Handler(loggedRouter)
	log.Fatal(http.ListenAndServe(p, handler))
}
