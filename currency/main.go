package main

import (
	"net"
	"os"

	"github.com/sgbaotran/Nascita-coffee-shop/currency/data"
	protos "github.com/sgbaotran/Nascita-coffee-shop/currency/protos/currency"
	"github.com/sgbaotran/Nascita-coffee-shop/currency/server"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewExchangeRates(log)
	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	cs := server.NewCurrencyServer(rates, log)

	// register the currency server
	protos.RegisterCurrencyServer(gs, cs)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Something goes wrong")
		os.Exit(1)
	}

	log.Info("[STARTING]: Currency-API Server Ready")
	// listen for requests
	gs.Serve(l)
}
