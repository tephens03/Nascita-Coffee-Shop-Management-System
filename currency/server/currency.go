package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/sgbaotran/Nascita-coffee-shop/currency/data"
	protos "github.com/sgbaotran/Nascita-coffee-shop/currency/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	rates         *data.ExchangeRate
	log           hclog.Logger
	subscriptions map[protos.Currency_SubscribeRateServer][]*protos.RateRequest
}

// NewCurrency creates a new Currency server
func NewCurrencyServer(rates *data.ExchangeRate, log hclog.Logger) *Currency {
	c := &Currency{log: log, rates: rates, subscriptions: map[protos.Currency_SubscribeRateServer][]*protos.RateRequest{}}
	go c.handleUpdates()
	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("New rate update")

		for key, rr := range c.subscriptions {

			for _, v := range rr {
				new_rate, err := c.rates.GetRates(v.GetBase().String(), v.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get rate: ", v.GetBase())
				}

				err = key.Send(&protos.RateResponse{Base: v.GetBase(), Destination: v.GetDestination(), Rate: new_rate})

			}

		}
	}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, req *protos.RateRequest) (*protos.RateResponse, error) {
	if req.Destination == req.Base {
		err := status.Newf(codes.InvalidArgument, "Base(%s) and Destination(%s) cannot be the same", req.Base, req.Destination)

		err, wde := err.WithDetails(req)
		if wde != nil {
			return nil, wde
		}
		return nil, err.Err()
	}

	c.log.Info("Handle Request for GetRate, base: ", req.GetBase(), ", destination: ", req.Destination)

	rate, err := c.rates.GetRates(req.Base.String(), req.Destination.String())

	if err != nil {
		return nil, err
	}
	return &protos.RateResponse{Base: req.Base, Destination: req.Destination, Rate: rate}, nil
}

func (c *Currency) SubscribeRate(src protos.Currency_SubscribeRateServer) error {

	for {
		rr, err := src.Recv()

		if err == io.EOF {
			c.log.Info("Client disconnected")
			delete(c.subscriptions, src)
			break
		}
		if err != nil {
			c.log.Error("Unable to read from client")
			return err
		}
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}
		rrs = append(rrs, rr)
		c.log.Info("New Rate: ", rr)
		c.subscriptions[src] = rrs
	}

	return nil
}
