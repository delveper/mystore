package rest

import (
	"net/http"

	"github.com/delveper/mystore/lib/lgr"
)

const productPath = "/products"

type Product struct {
	logger *lgr.Logger
}

func NewProduct(logger *lgr.Logger) Product {
	return Product{logger: logger}
}

func (p *Product) Create(rw http.ResponseWriter, req *http.Request) {
	p.logger.Info("bla bla bla")
	respondErr(rw, req, http.StatusInternalServerError, ErrNotImplemented)
}

func (p *Product) Read(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, ErrNotImplemented)
}

func (p *Product) ReadAll(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, ErrNotImplemented)
}

func (p *Product) Update(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, ErrNotImplemented)
}

func (p *Product) Delete(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, ErrNotImplemented)
}

func (p *Product) HandleEndpoint(mux *Mux) {
	hdl := func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {

		case http.MethodPost:
			p.Create(rw, req)

		case http.MethodGet:
			if _, ok := getPathID(req.URL.Path); ok {
				p.Read(rw, req) // it's possible to send parsed id via ctx
			} else {
				p.ReadAll(rw, req)
			}

		case http.MethodPut:
			p.Update(rw, req)

		case http.MethodDelete:
			p.Delete(rw, req)

		default:
			respondHTTPErr(rw, req, http.StatusNotFound)
		}
	}

	mux.HandleFunc(productPath, hdl)
}
