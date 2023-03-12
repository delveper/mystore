package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func decodeBody(req *http.Request, val any) error {
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(val); err != nil {
		return fmt.Errorf("decoding request body: %w", err)
	}

	return nil
}

func encodeBody(rw http.ResponseWriter, val any) error {
	if err := json.NewEncoder(rw).Encode(val); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}

func respond(rw http.ResponseWriter, req *http.Request, code int, data any) {
	rw.WriteHeader(code)

	if data != nil {
		if err := encodeBody(rw, data); err != nil {
			respondErr(rw, req, http.StatusInternalServerError)
		}
	}
}

func respondErr(rw http.ResponseWriter, req *http.Request, code int, args ...any) {
	respond(rw, req, code, map[string]any{
		"error": map[string]any{
			"message": fmt.Sprint(args...),
		},
	})
}

func respondHTTPErr(rw http.ResponseWriter, req *http.Request, code int) {
	respondErr(rw, req, code, http.StatusText(code))
}
