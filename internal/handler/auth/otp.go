package auth

import (
	"net/http"
)

func (a *authHandler) generateOTP(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) verifyOTP(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) validateOTP(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) disableOTP(w http.ResponseWriter, r *http.Request) {}
