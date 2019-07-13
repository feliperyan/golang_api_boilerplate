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
		fmt.Println(e)
	}

	if compose := os.Getenv("indocker"); compose == "dockercompose" {
		fmt.Println("Sleeping for 5 to wait for DB")
		time.Sleep(time.Second * 5)
	}

	if needsAuth := os.Getenv("NEEDS_AUTH"); needsAuth == "yes" {
		fmt.Println("Readying DB")
		readyDB()
	}

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.Use(jwtAuthentication)

	router.HandleFunc("/api/user/new", createAccount).Methods("POST")
	router.HandleFunc("/api/user/login", authenticate).Methods("POST")
	router.HandleFunc("/api/quote", quoteResponse).Methods("GET")
	router.HandleFunc("/api/quote/fr", quoteResponseFrench).Methods("GET")

	p := fmt.Sprintf(":%v", port)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Fatal(http.ListenAndServe(p, loggedRouter))
}
