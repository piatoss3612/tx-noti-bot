package auth

import "net/http"

func (a *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// read payload

	// create user data

	// insert user data to DB

	// write response
}

func (a *authHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// read payload

	// delete user data from DB by user id

	// write response
}
