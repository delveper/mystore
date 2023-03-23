package rest

import (
	"context"
	"errors"
	"fmt"
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

// Product handles product endpoint.
type Product struct {
	logic  ProductLogic
	logger lgr.Logger
}

// ProductDTO represents entities.Product model specific to transport layer.
type ProductDTO struct {
	ID          int    `json:"id"`
	MerchantID  int    `json:"merchant_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func NewProduct(logic ProductLogic, logger lgr.Logger) Product {
	return Product{
		logic:  logic,
		logger: logger,
	}
}

func withProductEntity(ctx context.Context, prod *entities.Product) context.Context {
	return context.WithValue(ctx, productKey, prod)
}

func getProductEntity(ctx context.Context) (*entities.Product, bool) {
	prod, ok := ctx.Value(productKey).(*entities.Product)

	return prod, ok
}

func convertProductsToDTO(prods ...entities.Product) []ProductDTO {
	res := make([]ProductDTO, 0, len(prods))

	for _, p := range prods {
		prod := ProductDTO{
			ID:          p.ID,
			MerchantID:  p.MerchantID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}

		res = append(res, prod)
	}

	return res
}

func (p Product) Create(rw http.ResponseWriter, req *http.Request) {
	var prod entities.Product
	if err := decodeBody(req, &prod); err != nil {
		respond(rw, http.StatusBadRequest, exceptions.ErrInvalidData)
		p.logger.Errorw("Failed to decode product from request body.", "error", err)

		return
	}

	defer func() {
		if err := req.Body.Close(); err != nil {
			p.logger.Warnf("Failed closing request %+v: %+v", req, err)
		}
	}()

	if err := prod.OK(); err != nil {
		respond(rw, http.StatusBadRequest, err)
		p.logger.Debugw("Failed validating product.", "error", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	id, err := p.logic.Add(ctx, prod)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			respond(rw, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrRecordExists):
			respond(rw, http.StatusConflict, exceptions.ErrRecordExists)
		case errors.Is(err, exceptions.ErrMerchantNotFound):
			respond(rw, http.StatusConflict, exceptions.ErrMerchantNotFound)
		default:
			respond(rw, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		p.logger.Errorw("Failed to add product.", "error", err)

		return
	}

	data := Response{Message: "Product created.", Details: fmt.Sprintf("id: %d", id)}

	respond(rw, http.StatusCreated, data)
	p.logger.Debugw(data.Message, "id", id)
}

func (p Product) Read(rw http.ResponseWriter, req *http.Request) {
	prod, ok := getProductEntity(req.Context())
	if !ok {
		respond(rw, http.StatusInternalServerError, exceptions.ErrInvalidData)
		p.logger.Errorw("Product not found in context.", "error", exceptions.ErrEmptyContext)

		return
	}

	res := convertProductsToDTO(*prod)

	respond(rw, http.StatusOK, res[0])
	p.logger.Debugw("Product Read successfully.")
}

func (p Product) ReadAll(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	prods, err := p.logic.FindMany(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			respond(rw, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrNotFound):
			respond(rw, http.StatusNotFound, exceptions.ErrNotFound)
		default:
			respond(rw, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		p.logger.Errorw("Failed find products.", "error", err)

		return
	}

	res := struct {
		Products []ProductDTO `json:"products"`
	}{convertProductsToDTO(prods...)}

	respond(rw, http.StatusOK, res)
	p.logger.Debugw("Products Read successfully.")
}

func (p Product) Update(rw http.ResponseWriter, req *http.Request) {
	prod, ok := getProductEntity(req.Context())
	if !ok {
		respond(rw, http.StatusInternalServerError, exceptions.ErrNotFound)
		p.logger.Errorw("Product not found", "error", exceptions.ErrEmptyContext)

		return
	}

	var update ProductDTO
	if err := decodeBody(req, &update); err != nil {
		respond(rw, http.StatusBadRequest, err)

		return
	}

	defer func() {
		if err := req.Body.Close(); err != nil {
			p.logger.Warnf("Failed closing request %+v: %+v", req, err)
		}
	}()

	if err := prod.OK(); err != nil {
		respond(rw, http.StatusBadRequest, err)
		p.logger.Debugw("Failed validating product.", "error", err)

		return
	}

	{ // modify existing product
		prod.Name = update.Name
		prod.Description = update.Description
		prod.Price = update.Price
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	if err := p.logic.Modify(ctx, *prod); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			respond(rw, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrNotFound):
			respond(rw, http.StatusNotFound, exceptions.ErrNotFound)
		default:
			respond(rw, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		p.logger.Errorw("Failed to Update product.", "error", err)

		return
	}

	respond(rw, http.StatusOK, Response{Message: "Product updated successfully"})
	p.logger.Debugw("Product updated successfully.", "product", prod)
}

func (p Product) Delete(rw http.ResponseWriter, req *http.Request) {
	prod, ok := getProductEntity(req.Context())
	if !ok {
		respond(rw, http.StatusInternalServerError, exceptions.ErrNotFound)
		p.logger.Errorw("Product not found", "error", exceptions.ErrEmptyContext)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	id, err := p.logic.Remove(ctx, *prod)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			respond(rw, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrNotFound):
			respond(rw, http.StatusNotFound, exceptions.ErrNotFound)
		default:
			respond(rw, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		p.logger.Errorw("Failed to Delete product.", "error", err)
	}

	p.logger.Debugw("Product deleted", "id", id)
}

// ServeHTTP implements http.Handler and handles Product's endpoint.
func (p Product) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Handle striped /products/ path
	if req.URL.Path == "" {
		switch req.Method {
		case http.MethodPost:
			p.Create(rw, req)
		case http.MethodGet:
			p.ReadAll(rw, req)
		default:
			respond(rw, http.StatusMethodNotAllowed, nil)
		}

		return
	}

	// Get ID from path
	id, err := strconv.Atoi(req.URL.Path)
	if err != nil {
		respond(rw, http.StatusBadRequest, nil)

		return
	}

	// Get Product by ID
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	prod, err := p.logic.Find(ctx, entities.Product{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			respond(rw, http.StatusGatewayTimeout, exceptions.ErrDeadline)
		case errors.Is(err, exceptions.ErrNotFound):
			respond(rw, http.StatusNotFound, nil)
		default:
			respond(rw, http.StatusInternalServerError, exceptions.ErrUnexpected)
		}

		p.logger.Errorw("Failed find product.", "error", err)

		return
	}
	// Send Product via context
	req = req.WithContext(withProductEntity(req.Context(), prod))

	// Handle path with ID
	switch req.Method {
	case http.MethodGet:
		p.Read(rw, req)
	case http.MethodPut:
		p.Update(rw, req)
	case http.MethodDelete:
		p.Delete(rw, req)
	default:
		respond(rw, http.StatusMethodNotAllowed, nil)
	}
}

// Route registers Product's endpoint pattern and strips root path.
func (p Product) Route(mux *http.ServeMux) {
	mux.Handle(productPath, http.StripPrefix(productPath, p))
}
