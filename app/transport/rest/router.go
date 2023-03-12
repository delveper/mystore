package rest

import "net/http"

type Mux struct{ *http.ServeMux }

func NewRouter(handlers ...func(*Mux)) *Mux {
	mux := &Mux{http.NewServeMux()}

	for i := range handlers {
		handlers[i](mux)
	}

	return mux
}
