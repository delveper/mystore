package rest

import (
	"net/http"
)

func ChainMiddlewares(hdl http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- { // LIFO order
		hdl = middlewares[i](hdl)
	}

	return hdl
}

func WithCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		if req.Method == "OPTIONS" {
			respondHTTPErr(rw, req, http.StatusNoContent)
			return
		}

		h(rw, req)
	}
}
