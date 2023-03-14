package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/delveper/mystore/lib/lgr"
)

type contextKey int

const contextKeyID contextKey = iota + 1

const productPath = "/products/"

type Product struct {
	logger lgr.Logger
}

func NewProduct(logger lgr.Logger) Product {
	return Product{logger: logger}
}

func (p *Product) create(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, errors.New("create"))
}

func (p *Product) read(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, errors.New("read"))
}

func (p *Product) readAll(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, errors.New("read all"))
}

func (p *Product) update(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, errors.New("update"))
}

func (p *Product) delete(rw http.ResponseWriter, req *http.Request) {
	respondErr(rw, req, http.StatusInternalServerError, errors.New("delete"))
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		switch req.Method {
		case http.MethodPost:
			p.create(rw, req)
		case http.MethodGet:
			p.readAll(rw, req)
		default:
			respondErr(rw, req, http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := strconv.Atoi(req.URL.Path)
	if err != nil {
		respondErr(rw, req, http.StatusNotFound)
		return
	}

	ctx := context.WithValue(req.Context(), contextKeyID, id)
	req = req.WithContext(ctx)

	switch req.Method {
	case http.MethodGet:
		p.read(rw, req)
	case http.MethodPut:
		p.update(rw, req)
	case http.MethodDelete:
		p.delete(rw, req)
	default:
		respondErr(rw, req, http.StatusMethodNotAllowed)
	}
}

func (p *Product) Route(mux *http.ServeMux) {
	mux.Handle(productPath, http.StripPrefix(productPath, p))
}
