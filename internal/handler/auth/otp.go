package auth

import "net/http"

func (a *authHandler) EnableOTP(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) DisableOTP(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) ValidateOTP(w http.ResponseWriter, r *http.Request) {}
