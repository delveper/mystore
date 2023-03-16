package rest

import "net/http"

// NewMux creates a new instance of *http.ServeMux and returns it
// after iteration over the input functions that calls newly created instance *http.ServeMux.
// This approach allows to send endpoint handlers in form of variadic functions that will address routing problem.
func NewMux(routes ...func(mux *http.ServeMux)) *http.ServeMux {
	mux := http.NewServeMux()

	for i := range routes {
		routes[i](mux)
	}

	return mux
}
