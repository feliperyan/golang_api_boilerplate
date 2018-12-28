package main

import (
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	userAccount := &UserAccount{}
	err := json.NewDecoder(r.Body).Decode(userAccount) //decode the request body into struct and failed if any error occur
	if err != nil {
		Respond(w, Message(false, "Invalid request"))
		return
	}

	resp := userAccount.Create() //Create account
	Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	userAccount := &UserAccount{}
	err := json.NewDecoder(r.Body).Decode(userAccount) //decode the request body into struct and failed if any error occur
	if err != nil {
		Respond(w, Message(false, "Invalid request"))
		return
	}

	resp := Login(userAccount.Email, userAccount.Password)
	Respond(w, resp)
}

func QuoteResponse(w http.ResponseWriter, r *http.Request) {
	resp := Message(true, "Success")
	resp["data"] = GetRandomQuote()
	Respond(w, resp)
}
