package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sgbaotran/Nascita-coffee-shop/product-api/data"
)

// ErrInvalidProductPath is an error message when the product path is not valid
func ValidateProductMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var prod data.Product

		err := data.FromJSON(prod, r.Body)

		if err != nil {
			http.Error(rw, "POST: Something went wrong (failed to serialize json) ("+(err.Error())+" :(", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(rw, "POST: Something went wrong (validation fail) ("+(err.Error())+" :(", http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

		r = r.WithContext(ctx)

		fmt.Println("In Middleware: ")

		next.ServeHTTP(rw, r)
	})

}
