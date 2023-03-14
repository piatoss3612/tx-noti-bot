package auth

import (
	"errors"
	"net/http"

	"github.com/piatoss3612/tx-notification/internal/helpers"
	"github.com/piatoss3612/tx-notification/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

var ErrInvalidUserID = errors.New("invalid user id")

func (a *authHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readUserPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while read and verify payload", err, "uri", r.RequestURI)
		return
	}

	var user models.User

	user.ID = payload.ID
	user.Email = payload.Email
	user.DiscordID = payload.DiscordID

	err = a.repo.CreateUser(r.Context(), &user)
	if err != nil {
		errMsg := "unable to register"

		if mongo.IsDuplicateKeyError(err) {
			errMsg = "user already registered"
		}

		_ = helpers.ErrorJSON(w, http.StatusBadRequest, errMsg)
		slog.Error("error while creating user", err, "uri", r.RequestURI)
		return
	}

	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readUserPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while read and verify payload", err, "uri", r.RequestURI)
		return
	}

	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID, "uri", r.RequestURI)
		return
	}

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

func (a *authHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readUserPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while read and verify payload", err, "uri", r.RequestURI)
		return
	}

	user, err := a.repo.GetUserByID(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while retrieving user matched by id", err, "id", payload.ID, "uri", r.RequestURI)
		return
	}

	user.Email = payload.Email
	user.DiscordID = payload.DiscordID

	err = a.repo.UpdateUser(r.Context(), user)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "unable to generate otp")
		slog.Error("error while updating user", err, "uri", r.RequestURI)
		return
	}

	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	_ = helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	payload, err := a.readUserPayloadAndVerify(w, r)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, err.Error())
		slog.Error("error while read and verify payload", err, "uri", r.RequestURI)
		return
	}

	err = a.repo.DeleteUser(r.Context(), payload.ID)
	if err != nil {
		_ = helpers.ErrorJSON(w, http.StatusBadRequest, "invalid user id or non-existing user")
		slog.Error("error while deleting user matched by id", err, "id", payload.ID, "uri", r.RequestURI)
		return
	}

	var resp models.CommonResponse

	resp.StatusCode = http.StatusOK

	helpers.WriteJSON(w, http.StatusOK, resp)
}

func (a *authHandler) readUserPayloadAndVerify(w http.ResponseWriter, r *http.Request) (*models.UserPayload, error) {
	var payload models.UserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		return nil, err
	}

	if !helpers.IsValidEthereumAddress(payload.ID) {
		return nil, ErrInvalidUserID
	}

	return &payload, nil
}
