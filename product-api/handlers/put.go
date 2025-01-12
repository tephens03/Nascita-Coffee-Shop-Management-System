package handlers

import (
	"net/http"
	"strconv"

	"github.com/sgbaotran/Nascita-coffee-shop/product-api/data"

	"github.com/gorilla/mux"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "PUT: Something went wrong (cannot convert ID) ("+err.Error()+") :(", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = p.productDB.UpdateProduct(id, &prod)
	if err != nil {
		http.Error(rw, "PUT: Something went wrong ("+err.Error()+") :(", http.StatusBadRequest)
		return
	}
	p.l.Info("Updated products: ", prod, err)

}
