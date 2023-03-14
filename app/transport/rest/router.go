package rest

import "net/http"

func NewMux(routes ...func(*http.ServeMux)) *http.ServeMux {
	mux := http.NewServeMux()

	for i := range routes {
		routes[i](mux)
	}

	return mux
}
