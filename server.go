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

	//ReadyDB()

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.Use(JwtAuthentication)

	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/quote", QuoteResponse).Methods("GET")

	p := fmt.Sprintf(":%v", port)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowCredentials: true,
	// 	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	// 	Debug:            true,
	// })
	// handler := c.Handler(loggedRouter)

	// log.Fatal(http.ListenAndServe(p, handler))

	log.Fatal(http.ListenAndServe(p, loggedRouter))
}
