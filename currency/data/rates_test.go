package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func Test(t *testing.T) {
	l := hclog.Default()
	ex, err := NewExchangeRates(l)

	if err != nil {
		t.Fatal(err)
	}

	err = ex.getRates()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v", ex.rates)

}
