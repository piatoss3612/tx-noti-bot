package auth

import "net/http"

func (a *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) LoginUser(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}
