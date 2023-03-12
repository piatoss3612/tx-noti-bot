package auth

import (
	"net/http"
	"regexp"

	"github.com/piatoss3612/tx-noti-bot/internal/helpers"
	"github.com/piatoss3612/tx-noti-bot/internal/models"
	"golang.org/x/exp/slog"
)

func (a *authHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.UserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid payload")
		slog.Error("error while reading json", err)
		return
	}

	// check if id is valid ethereum address
	if !isValidAddress(payload.ID) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// create new user
	var user models.User

	user.ID = payload.ID
	user.Email = payload.Email
	user.DiscordID = payload.DiscordID

	// insert user into db
	err = a.repo.CreateUser(r.Context(), &user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "")
		slog.Error("error while creating new user", err)
		return
	}

	// send response
	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.UserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid payload")
		slog.Error("error while reading json", err)
		return
	}

	// check if id is valid ethereum address
	if !isValidAddress(payload.ID) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// get user by id
	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID)
		return
	}

	// send response
	var resp models.UserResponse

	resp.StatusCode = http.StatusOK
	resp.User.ID = user.ID
	resp.User.Email = user.Email
	resp.User.DiscordID = user.DiscordID
	resp.User.OtpEnabled = user.OtpEnabled
	resp.User.OtpVerified = user.OtpVerified
	resp.User.CreatedAt = user.CreatedAt

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.UserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid payload")
		slog.Error("error while reading json", err)
		return
	}

	// check if id is valid ethereum address
	if !isValidAddress(payload.ID) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// delete user by id
	err = a.repo.DeleteUser(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while deleting user matched by id", err, "id", payload.ID)
		return
	}

	// send response
	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	helpers.WriteJSON(w, http.StatusOK, resp)
}

func isValidAddress(s string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(s)
}
