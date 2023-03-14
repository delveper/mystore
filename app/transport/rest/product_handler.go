package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/delveper/mystore/app/entities"
	"github.com/delveper/mystore/app/exceptions"
	"github.com/delveper/mystore/lib/lgr"
)

type contextKey int

const productKey contextKey = iota + 1

const productPath = "/products/"

const queryTimeout = 3 * time.Second

type Product struct {
	logic  ProductLogic
	logger lgr.Logger
}

func NewProduct(logic ProductLogic, logger lgr.Logger) Product {
	return Product{
		logic:  logic,
		logger: logger,
	}
}

func withProduct(ctx context.Context, prod *entities.Product) context.Context {
	return context.WithValue(ctx, productKey, prod)
}

func getProduct(ctx context.Context) (*entities.Product, bool) {
	prod, ok := ctx.Value(productKey).(*entities.Product)

	return prod, ok
}

func (p *Product) create(rw http.ResponseWriter, req *http.Request) {
	respond(rw, http.StatusInternalServerError, errors.New("create"))
}

func (p *Product) read(rw http.ResponseWriter, req *http.Request) {
	respond(rw, http.StatusInternalServerError, errors.New("read"))
}

func (p *Product) readAll(rw http.ResponseWriter, req *http.Request) {
	respond(rw, http.StatusInternalServerError, errors.New("read all"))
}

func (p *Product) update(rw http.ResponseWriter, req *http.Request) {
	respond(rw, http.StatusInternalServerError, errors.New("update"))
}

func (p *Product) delete(rw http.ResponseWriter, req *http.Request) {
	respond(rw, http.StatusInternalServerError, errors.New("delete"))
}

// ServeHTTP implements http.Handler and handles Product's endpoint.
func (p *Product) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Handle striped path
	if req.URL.Path == "" {
		switch req.Method {
		case http.MethodPost:
			p.create(rw, req)
		case http.MethodGet:
			p.readAll(rw, req)
		default:
			respond(rw, http.StatusMethodNotAllowed, nil)
		}

		return
	}

	// Get ID from path
	id, err := strconv.Atoi(req.URL.Path)
	if err != nil {
		respond(rw, http.StatusNotFound, nil)
		return
	}

	// Get Product by ID
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	prod, err := p.logic.Find(ctx, entities.Product{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, exceptions.ErrDeadline):
			respond(rw, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		default:
			respond(rw, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}
		p.logger.Errorw("Failed find product.", "error", err)

		return
	}
	// Send Product via context
	req = req.WithContext(withProduct(req.Context(), prod))

	// Handle path with ID
	switch req.Method {
	case http.MethodGet:
		p.read(rw, req)
	case http.MethodPut:
		p.update(rw, req)
	case http.MethodDelete:
		p.delete(rw, req)
	default:
		respond(rw, http.StatusMethodNotAllowed, nil)
	}
}

// Route registers Product's endpoint pattern and strips root path.
func (p *Product) Route(mux *http.ServeMux) {
	mux.Handle(productPath, http.StripPrefix(productPath, p))
}
