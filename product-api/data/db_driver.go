package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	protos "github.com/sgbaotran/Nascita-coffee-shop/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
)

type ProductsDB struct {
	cc     protos.CurrencyClient
	client protos.Currency_SubscribeRateClient
	l      hclog.Logger
	rates  map[string]float64
}

func NewProductsSB(cc protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{cc, nil, l, map[string]float64{}}
	go pb.handleUpdates()
	return pb
}

// GetProducts returns the list of products.
func (pdb *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return products, nil
	}
	rate, err := pdb.GetExchangeRate(currency)

	if err != nil {
		pdb.l.Error("[ERROR]:[PRODUCT_DB]: Unable to get exchange rate ", err)
		return nil, err
	}
	new_prods := Products{}
	for _, product := range products {
		np := *product
		np.Price *= rate
		new_prods = append(new_prods, &np)
	}

	return new_prods, err
}

func (pdb *ProductsDB) handleUpdates() {
	client, err := pdb.cc.SubscribeRate(context.Background())

	if err != nil {
		pdb.l.Error("Unable to subcribe")
	}
	pdb.client = client
	for {
		rr, err := client.Recv()
		if err != nil {
			pdb.l.Error("Unable to receive message 	")
			return
		}
		pdb.rates[rr.Destination.String()] = rr.Rate

	}
}

// GetProduct returns a product by its ID.
func (pdb *ProductsDB) GetProduct(id int, currency string) (*Product, error) {

	i := pdb.findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}
	prod := *products[i]

	if currency == "" {
		return &prod, nil
	}

	rate, err := pdb.GetExchangeRate(currency)

	if err != nil {
		pdb.l.Error("[ERROR]:[PRODUCT_DB]: Unable to get exchange rate ", err)
		return nil, err
	}

	prod.Price *= rate

	return &prod, nil
}

// AddProduct adds a new product to the database.
func (pdb *ProductsDB) AddProduct(p *Product) {
	p.ID = products[len(products)-1].ID + 1
	products = append(products, p)
}

// findProduct finds a product by its ID and returns it along with its index in the slice.
func (pdb *ProductsDB) findProduct(id int) (*Product, int, error) {
	for ind, prod := range products {
		if prod.ID == id {
			return prod, ind, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

// UpdateProduct updates an existing product in the database.
func (pdb *ProductsDB) UpdateProduct(id int, p *Product) error {
	old_prod, ind, err := pdb.findProduct(id)
	if err != nil {
		return err
	}
	p.ID = old_prod.ID
	products[ind] = p
	return nil
}

// DeleteProduct deletes a product from the database.
func (pdb *ProductsDB) DeleteProduct(id int) error {
	i := pdb.findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}
	products = append(products[:i], products[i+1])
	return nil
}

// findIndexByProductID finds the index of a product in the database.
// It returns -1 when no product can be found.
func (pdb *ProductsDB) findIndexByProductID(id int) int {
	for i, p := range products {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// ErrProductNotFound is an error returned when a product is not found.
var ErrProductNotFound = fmt.Errorf("Product not found")

// products stores a list of Product pointers.
var products = []*Product{
	&Product{
		ID:          1,
		Name:        "Laptop",
		Description: "Powerful laptop with high-performance specifications",
		Price:       999.99,
		SKU:         "ABC-ABC-ABC",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Smartphone",
		Description: "Feature-rich smartphone with a high-quality camera",
		Price:       599.99,
		SKU:         "ABC-ABC-ABC",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

func (pdb *ProductsDB) GetExchangeRate(destination string) (float64, error) {
	destination = strings.ToUpper(destination)

	if r, ok := pdb.rates[destination]; ok {
		return r, nil
	}

	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_CAD),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := pdb.cc.GetRate(context.Background(), rr)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0]

			i, _ := md.(*protos.RateRequest)

			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("Unable to get rate, destination and base currencies cannot be the same Base: %s, Destination: %s", i.Base.String(), i.Destination.String())
			}

			return -1, fmt.Errorf("Unable to get rate, destination and base currencies cannot be the same Base: %s, Destination: %s", i.Base.String(), i.Destination.String())
		}
		return -1, err
	}

	pdb.rates[destination] = resp.Rate

	pdb.client.Send(rr)

	return resp.Rate, nil
}
