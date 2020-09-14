package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Respond marshals a value to JSON and sends it to the client.
func Respond(writer http.ResponseWriter, val interface{}, statusCode int) error {
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "marshalling value to json")
	}

	writer.Header().Set("content-type", "application/json; charset=utf-8")
	writer.WriteHeader(statusCode)

	if _, err := writer.Write(data); err != nil {
		return errors.Wrap(err, "writing to client")
	}

	return nil
}

// RespondError knows how to handle errors going out to the client.
func RespondError(writer http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		resp := ErrorResponse{
			Error: webErr.Err.Error(),
		}

		return Respond(writer, resp, webErr.Status)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(writer, resp, http.StatusInternalServerError)
}
