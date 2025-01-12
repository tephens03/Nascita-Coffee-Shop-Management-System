package data

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Products is a slice of Product pointers.
type Products []*Product

// Product represents a product with basic information.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// validateSKU checks if a SKU follows a specific format.
func validateSKU(fl validator.FieldLevel) bool {

	reg := regexp.MustCompile("[A-Z]+-[A-Z]+-[A-Z]+")
	matches := reg.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true
}

// Validate performs validation on a Product instance.
func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}
