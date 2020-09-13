package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Decode looks for a JSON document in the request body and unmarshalls it into val.
func Decode(request *http.Request, val interface{}) error {
	if err := json.NewDecoder(request.Body).Decode(val); err != nil {
		return errors.Wrap(err, "decoding request body")
	}
	return nil
}
