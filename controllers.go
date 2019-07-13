package main

import (
	"encoding/json"
	"net/http"
)

var createAccount = func(w http.ResponseWriter, r *http.Request) {
	someUser := &userAccount{}
	err := json.NewDecoder(r.Body).Decode(someUser) //decode the request body into struct and failed if any error occur
	if err != nil {
		respond(w, message(false, "Invalid request"))
		return
	}

	resp := someUser.Create() //Create account
	respond(w, resp)
}

var authenticate = func(w http.ResponseWriter, r *http.Request) {
	someUser := &userAccount{}
	err := json.NewDecoder(r.Body).Decode(someUser) //decode the request body into struct and failed if any error occur
	if err != nil {
		respond(w, message(false, "Invalid request"))
		return
	}

	resp := login(someUser.Email, someUser.Password)
	respond(w, resp)
}

func quoteResponse(w http.ResponseWriter, r *http.Request) {
	resp := message(true, "Success")
	resp["data"] = getRandomQuote("English")
	respond(w, resp)
}

func quoteResponseFrench(w http.ResponseWriter, r *http.Request) {
	resp := message(true, "Success")
	resp["data"] = getRandomQuote("French")
	respond(w, resp)
}
