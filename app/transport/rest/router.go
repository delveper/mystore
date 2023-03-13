package rest

import "net/http"

func NewMux(hds ...func(mux *http.ServeMux)) *http.ServeMux {
	mux := http.NewServeMux()

	for i := range hds {
		hds[i](mux)
	}

	return mux
}
