package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/piatoss3612/tx-noti-bot/internal/models"
)

const MaxRequestBodyBytes = 1 << 20

var (
	ErrRequestBodyWithMultipleJSON = errors.New("request body must contain a single JSON value")
)

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	defer func() { _ = r.Body.Close() }()

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodyBytes)

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return ErrRequestBodyWithMultipleJSON
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any, headers ...http.Header) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) != 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(body)

	return err
}

func ErrorJSON(w http.ResponseWriter, statusCode int, errMsg string) error {
	var resp models.CommonResponse

	resp.StatusCode = statusCode
	resp.Error = errMsg

	return WriteJSON(w, statusCode, resp)
}
