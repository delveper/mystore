package rest

import (
	"net/http"

	"github.com/delveper/mystore/lib/lgr"
)

func ChainMiddlewares(hdl http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- { // LIFO
		hdl = middlewares[i](hdl)
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
		rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		if req.Method == http.MethodOptions {
			respond(rw, http.StatusNoContent, nil)

			return
		}

		hdl.ServeHTTP(rw, req)
	})
}

func WithAuth(hdl http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		userName, password, ok := req.BasicAuth()
		if !(ok && isAuth(userName, password)) {
			rw.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			respond(rw, http.StatusUnauthorized, nil)

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

// WithoutPanic recovers in case panic, but we won't panic.
func WithoutPanic(logger lgr.Logger) func(http.Handler) http.Handler {
	return func(hdl http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					respond(rw, http.StatusInternalServerError, nil)
					logger.Errorw("Recovered from panic.", "error", err)
				}
			}()

			hdl.ServeHTTP(rw, req)
		})
	}
}
