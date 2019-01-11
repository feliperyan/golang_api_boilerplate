package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
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

	if needsAuth := os.Getenv("NEEDS_AUTH"); needsAuth == "yes" {
		fmt.Println("Readying DB")
		ReadyDB()
	}

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.Use(JwtAuthentication)

	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/quote", QuoteResponse).Methods("GET")

	p := fmt.Sprintf(":%v", port)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(p, loggedRouter))
}
