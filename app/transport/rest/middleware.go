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

func WithJSON(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		hdl(rw, req)
	}
}

func WithCORS(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		if req.Method == "OPTIONS" {
			respondHTTPErr(rw, req, http.StatusNoContent)
			return
		}

		hdl(rw, req)
	}
}

func WithAuth(hdl http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok {
			rw.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			rw.WriteHeader(http.StatusUnauthorized)

			return
		}

		_, _ = username, password

		if false {
			rw.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			rw.WriteHeader(http.StatusUnauthorized)

			return
		}

		hdl(rw, req)
	}
}
