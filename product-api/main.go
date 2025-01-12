package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	protos "github.com/sgbaotran/Nascita-coffee-shop/currency/protos/currency"
	"github.com/sgbaotran/Nascita-coffee-shop/product-api/data"
	"github.com/sgbaotran/Nascita-coffee-shop/product-api/handlers"
	"google.golang.org/grpc"
)

func main() {
	// _ := log.New(os.Stdout, "REST-API ", log.LstdFlags)
	log := hclog.Default()
	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	/* mux := http.NewServeMux()
	mux.Handle("/", handlers.NewProduct(l))*/

	cc := protos.NewCurrencyClient(conn)

	pdb := data.NewProductsSB(cc, log)

	ph := handlers.NewProducts(cc, log, pdb)

	r := mux.NewRouter()

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.HandleFunc("/", ph.GetProducts).Queries("currency", "{[a-zA-Z]{3}}")

	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetProduct)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetProduct).Queries("currency", "{[a-zA-Z]{3}}")

	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(handlers.ValidateProductMiddleWare)

	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(handlers.ValidateProductMiddleWare)

	deleteRouter := r.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)
	deleteRouter.Use(handlers.ValidateProductMiddleWare)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	server := &http.Server{
		Addr:         ":3030",
		Handler:      r,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Error("Unable to start server", "error", err)
		}
	}()
	log.Info("[STARTING]: Product-API Server Ready")

	signalChan := make(chan os.Signal)

	signal.Notify(signalChan, os.Interrupt)

	signal.Notify(signalChan, os.Kill)

	sig := <-signalChan
	log.Info("Somebody turned off", sig)

}
