package auth

import (
	"net/http"

	"github.com/piatoss3612/tx-notification/internal/helpers"
	"github.com/piatoss3612/tx-notification/internal/models"
	"github.com/pquerna/otp/totp"
	"golang.org/x/exp/slog"
)

func (a *authHandler) generateOTP(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readOtpPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while reading and verifying payload", err, "uri", r.RequestURI)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "piatoss.tech",
		AccountName: payload.ID,
	})
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while generating otp key", err, "uri", r.RequestURI)
		return
	}

	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID, "uri", r.RequestURI)
		return
	}

	user.OtpVerified = false
	user.OtpSecret = key.Secret()
	user.OtpUrl = key.URL()

	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "unable to generate otp")
		slog.Error("error while updating user", err, "uri", r.RequestURI)
		return
	}

	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpSecret = key.Secret()
	resp.Otp.OtpUrl = key.URL()

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) verifyOTP(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readOtpPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while reading and verifying payload", err, "uri", r.RequestURI)
		return
	}

	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID, "uri", r.RequestURI)
		return
	}

	if !totp.Validate(payload.Token, user.OtpSecret) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid token or non-existing user")
		return
	}

	user.OtpEnabled = true
	user.OtpVerified = true

	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "unable to verify otp")
		slog.Error("error while updating user", err, "uri", r.RequestURI)
		return
	}

	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpEnabled = true
	resp.Otp.OtpVerified = true

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) validateOTP(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readOtpPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while reading and verifying payload", err, "uri", r.RequestURI)
		return
	}

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

	if !totp.Validate(payload.Token, user.OtpSecret) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid token or non-existing user")
		return
	}

	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpValid = true

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) disableOTP(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readOtpPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while reading and verifying payload", err, "uri", r.RequestURI)
		return
	}

	if !helpers.IsValidEthereumAddress(payload.ID) {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID, "uri", r.RequestURI)
		return
	}

	user.OtpEnabled = false
	user.OtpVerified = false

	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "non-existing user")
		slog.Error("error while updating user", err, "uri", r.RequestURI)
		return
	}

	var resp models.OtpResponse

	resp.StatusCode = http.StatusOK
	resp.Otp.OtpEnabled = false

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) readOtpPayloadAndVerify(w http.ResponseWriter, r *http.Request) (*models.OtpPayload, error) {
	var payload models.OtpPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		return nil, err
	}

	if !helpers.IsValidEthereumAddress(payload.ID) {
		return nil, ErrInvalidUserID
	}

	return &payload, nil
}
