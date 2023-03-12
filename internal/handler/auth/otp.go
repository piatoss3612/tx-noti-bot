package auth

import (
	"net/http"

	"github.com/piatoss3612/tx-noti-bot/internal/helpers"
	"github.com/piatoss3612/tx-noti-bot/internal/models"
	"github.com/pquerna/otp/totp"
	"golang.org/x/exp/slog"
)

func (a *authHandler) generateOTP(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.OtpPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while reading json", err)
		return
	}

	// check if user id valid ethereum address
	if !isValidAddress(payload.ID) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	// generate otp key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "piatoss3612",
		AccountName: payload.ID,
	})
	if err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while generating otp key", err)
		return
	}

	// get user matched by id
	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID)
		return
	}

	// update user state
	user.OtpEnabled = false
	user.OtpVerified = false
	user.OtpSecret = key.Secret()
	user.OtpUrl = key.URL()

	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "unable to generate otp")
		slog.Error("error while updating user", err)
		return
	}

	// send response
	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpSecret = key.Secret()
	resp.Otp.OtpUrl = key.URL()

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) verifyOTP(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.OtpPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
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

	// validate received token
	isValid := totp.Validate(payload.Token, user.OtpSecret)
	if !isValid {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid token or non-existing user")
		return
	}

	user.OtpEnabled = true
	user.OtpVerified = true

	// update user
	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "unable to verify otp")
		slog.Error("error while updating user", err)
		return
	}

	// send response
	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpEnabled = true
	resp.Otp.OtpVerified = true

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) validateOTP(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.OtpPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
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

	if !user.OtpEnabled || !user.OtpVerified {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "please reconfigure MFA")
		return
	}

	// verify token
	isValid := totp.Validate(payload.Token, user.OtpSecret)
	if !isValid {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid token or non-existing user")
		return
	}

	// send response
	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpValid = true

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) disableOTP(w http.ResponseWriter, r *http.Request) {
	// read payload
	var payload models.OtpPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
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

	user.OtpEnabled = false

	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "non-existing user")
		slog.Error("error while updating user", err)
		return
	}

	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpEnabled = false

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}
