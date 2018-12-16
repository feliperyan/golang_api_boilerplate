package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// our main function
func main() {

	router := mux.NewRouter()
	router.Use(JwtAuthentication)

	router.HandleFunc("/api/user/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/dummy", DummyResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
