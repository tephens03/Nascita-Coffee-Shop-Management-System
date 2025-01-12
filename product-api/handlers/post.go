package handlers

import (
	"net/http"

	"github.com/sgbaotran/Nascita-coffee-shop/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.productDB.AddProduct(&prod)

	p.l.Info("Added products: ", prod)

}
