package data

import (
	"testing"
)

func TestValidationFunctionality(t *testing.T) {
	p := Product{
		Name:  "Macbook Pro 14",
		Price: 22.0,
		SKU:   "a-a-a",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
