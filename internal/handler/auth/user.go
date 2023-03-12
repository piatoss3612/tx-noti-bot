package auth

import (
	"net/http"

	"github.com/piatoss3612/tx-noti-bot/internal/helpers"
	"github.com/piatoss3612/tx-noti-bot/internal/models"
)

func (a *authHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.UserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, "")
		return
	}

	// create user data
	var user models.User

	user.ID = payload.ID
	user.Email = payload.Email
	user.DiscordID = payload.DiscordID

	// insert user data to DB
	err = a.repo.CreateUser(r.Context(), &user)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, "")
		return
	}

	// write response
	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) loginUser(w http.ResponseWriter, r *http.Request) {}

func (a *authHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.UserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, "")
		return
	}

	// delete user data from DB by user id
	err = a.repo.DeleteUser(r.Context(), payload.ID)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, "")
		return
	}

	// write response
	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	helpers.WriteJSON(w, http.StatusOK, resp)
}
