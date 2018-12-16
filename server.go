package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// our main function
func main() {

	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.Use(JwtAuthentication)

	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/dummy", DummyResponse).Methods("GET")
	p := fmt.Sprintf(":%v", port)
	log.Fatal(http.ListenAndServe(p, router))
}
