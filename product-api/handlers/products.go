package handlers

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	protos "github.com/sgbaotran/Nascita-coffee-shop/currency/protos/currency"
	"github.com/sgbaotran/Nascita-coffee-shop/product-api/data"
)

// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body int
}

// Products handler for getting and updating products
type Products struct {
	l         hclog.Logger
	cc        protos.CurrencyClient
	productDB *data.ProductsDB
}

// NewProducts returns a new products handler with the given logger
func NewProducts(cc protos.CurrencyClient, l hclog.Logger, pdb *data.ProductsDB) *Products {
	return &Products{l, cc, pdb}
}

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

var ErrProductNotFound = fmt.Errorf("Product not found")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}
