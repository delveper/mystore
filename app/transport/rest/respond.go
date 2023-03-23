package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func decodeBody(req *http.Request, data any) error {
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return fmt.Errorf("decoding request body: %w", err)
	}

	return nil
}

func encodeBody(rw http.ResponseWriter, data any) error {
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}

func respond(rw http.ResponseWriter, code int, data any) {
	rw.WriteHeader(code)

	if err, ok := data.(error); ok {
		data = Response{Message: http.StatusText(code), Details: err.Error()}
	}

	if data == nil && code != http.StatusNoContent {
		data = Response{Message: http.StatusText(code)}
	}

	if data != nil {
		if err := encodeBody(rw, data); err != nil {
			respond(rw, http.StatusInternalServerError, err)
		}
	}
}
