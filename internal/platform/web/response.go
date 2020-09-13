package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)


// Respond marshalls a value to JSON and sends it to the client.
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
