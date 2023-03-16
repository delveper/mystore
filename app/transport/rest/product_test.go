package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/delveper/mystore/app/entities"
	"github.com/delveper/mystore/app/exceptions"
	"github.com/delveper/mystore/lib/lgr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestProduct_Create(t *testing.T) {
	logger := lgr.New()

	testsCases := map[string]struct {
		requestBody     string
		expCode         int
		expResponseBody string
		mockMethod      func(ctrl *gomock.Controller) ProductLogic
	}{
		"SUCCESS": {
			requestBody:     `{"merchant_id": 1, "name": "the thing", "description": "number one", "price": 1000}`,
			expCode:         http.StatusCreated,
			expResponseBody: `{"message": "Product created.", "details": "id: 1"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				mock := NewMockProductLogic(ctrl)

				mock.EXPECT().
					Add(gomock.Any(), entities.Product{MerchantID: 1, Name: "the thing", Description: "number one", Price: 1000}).
					Return(1, nil).Times(1)

				return mock
			},
		},
		"INVALID_REQUEST_BODY": {
			requestBody:     `invalid json`,
			expCode:         http.StatusBadRequest,
			expResponseBody: `{"message": "Bad Request", "details": "invalid data"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				return nil
			},
		},

		"SERVER_ERROR": {
			requestBody:     `{"merchant_id": 1, "name": "the thing", "description": "number one", "price": 1000}`,
			expCode:         http.StatusInternalServerError,
			expResponseBody: `{"message": "Internal Server Error", "details": "unexpected error"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				mock := NewMockProductLogic(ctrl)

				mock.EXPECT().
					Add(gomock.Any(), entities.Product{MerchantID: 1, Name: "the thing", Description: "number one", Price: 1000}).
					Return(0, exceptions.ErrUnexpected).Times(1)

				return mock
			},
		},

		"CONTEXT_TIMEOUT": {
			requestBody:     `{"merchant_id": 1, "name": "the thing", "description": "number one", "price": 1000}`,
			expCode:         http.StatusGatewayTimeout,
			expResponseBody: `{"message": "Gateway Timeout", "details": "deadline exceeded"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				// Mock a long-running request
				mock := NewMockProductLogic(ctrl)

				mock.EXPECT().
					Add(gomock.Any(), entities.Product{MerchantID: 1, Name: "the thing", Description: "number one", Price: 1000}).
					DoAndReturn(func(ctx context.Context, product entities.Product) (int, error) {
						time.Sleep(queryTimeout + time.Second)
						return 0, context.DeadlineExceeded
					}).Times(1)

				return mock
			},
		},

		"MERCHANT_VALIDATION_ERROR": {
			requestBody:     `{"merchant_id": -1, "name": "the thing", "description": "number one", "price": 1000}`,
			expCode:         http.StatusBadRequest,
			expResponseBody: `{"message": "Bad Request", "details": "merchant_id must be a positive integer" }`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				return nil
			},
		},
		"PRICE_VALIDATION_ERROR": {
			requestBody:     `{"merchant_id": 1, "name": "the thing", "description": "number one", "price": -1000}`,
			expCode:         http.StatusBadRequest,
			expResponseBody: `{"message":"Bad Request","details":"price must be a positive integer"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				return nil
			},
		},
		"NAME_VALIDATION_ERROR": {
			requestBody:     `{"merchant_id": 1, "name": "t", "description": "number one", "price": 1000}`,
			expCode:         http.StatusBadRequest,
			expResponseBody: `{"message":"Bad Request","details":"name must be between 2 and 255 characters long"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				return nil
			},
		},
		"CASCADE_VALIDATION_ERROR": {
			requestBody:     `{"merchant_id": -1, "name": "s", "description": "does not matter", "price": -1000}`,
			expCode:         http.StatusBadRequest,
			expResponseBody: `{"message":"Bad Request","details":"merchant_id must be a positive integer\nname must be between 2 and 255 characters long\nprice must be a positive integer"}`,
			mockMethod: func(ctrl *gomock.Controller) ProductLogic {
				return nil
			},
		},
	}

	for name, tt := range testsCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			p := NewProduct(tt.mockMethod(ctrl), logger)

			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, productPath, strings.NewReader(tt.requestBody))

			p.Create(resp, req)

			require.Equal(t, tt.expCode, resp.Code)
			require.JSONEq(t, tt.expResponseBody, resp.Body.String())
		})
	}
}
