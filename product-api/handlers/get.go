package handlers

import (
	"net/http"

	"github.com/sgbaotran/Nascita-coffee-shop/product-api/data"
)

// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Info("[PRODUCT_HANDLER]: Getting list of all products")

	rw.Header().Add("Content-Type", "application/json")

	cur := r.URL.Query().Get("currency")

	prods, err := p.productDB.GetProducts(cur)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prods, rw)

	if err != nil {
		http.Error(rw, "GET: Something went wrong :( "+err.Error(), http.StatusBadRequest)
		return
	}
	p.l.Info("[SUCCESS]:[PRODUCT_HANDLER]: Successfully getting list of all products")

}

// Return a product from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("[PRODUCT_HANDLER]: Getting product by ID")

	cur := r.URL.Query().Get("currency")

	id := getProductID(r)

	prod, err := p.productDB.GetProduct(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("[ERROR]:[PRODUCT_HANDLER]: Unable to fetch product", "error", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("[ERROR]:[PRODUCT_HANDLER]: Unable to fetching product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)

	// we should never be here but log the error just incase
	if err != nil {
		http.Error(rw, "GET: Something went wrong :( "+err.Error(), http.StatusBadRequest)
		return
	}

	p.l.Info("[SUCCESS]:[PRODUCT_HANDLER]: Successfully getting product")

}
