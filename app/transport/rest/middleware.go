package rest

import (
	"net/http"

	"github.com/delveper/mystore/lib/lgr"
)

func ChainMiddlewares(hdl http.Handler, mds ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mds) - 1; i >= 0; i-- { // LIFO order
		hdl = mds[i](hdl)
	}

	return hdl
}

func WithJSON(hdl http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		hdl.ServeHTTP(rw, req)
	})
}

func WithCORS(hdl http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		if req.Method == "OPTIONS" {
			respondHTTPErr(rw, req, http.StatusNoContent)
			return
		}

		hdl.ServeHTTP(rw, req)
	})
}

func WithAuth(hdl http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
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

		hdl.ServeHTTP(rw, req)
	})
}

// WithLogRequest logs every request and sends logger instance to further handler.
func WithLogRequest(logger lgr.Logger) func(http.Handler) http.Handler {
	return func(hdl http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Debugw("Request:",
				"method", req.Method,
				"uri", req.RequestURI,
				"user-agent", req.UserAgent(),
				"remote", req.RemoteAddr,
			)

			hdl.ServeHTTP(rw, req)
		})
	}
}
