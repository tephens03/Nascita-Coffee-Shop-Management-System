package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/sgbaotran/Nascita-coffee-shop/product-images/files"
	"github.com/sgbaotran/Nascita-coffee-shop/product-images/handlers"
)

// @Summary Get user by ID
// @Description Get user information by ID
func main() {

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString("debug"),
		},
	)

	// create a logger for the server from the default logger
	error_log := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// hc_l := hclog.Default()

	storage, err := files.NewLocal(1024*1000*5, "./imagestore")
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	file_handler := handlers.NewFile(l, storage)

	sm := mux.NewRouter()
	gh := handlers.GzipHandler{}

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", file_handler.UploadFileREST)
	postRouter.Use(gh.GzipMiddleware)
	//Since the data will be in the request, needs not to have the params
	// postRouter.HandleFunc("/", file_handler.UploadFileMultipart)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/images/", http.FileServer(http.Dir("./imagestore"))))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		ErrorLog:     error_log,
	}

	// start the server
	go func() {
		l.Info("Starting server", "bind_address", "9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal)

	signal.Notify(signalChan, os.Interrupt)

	signal.Notify(signalChan, os.Kill)

	sig := <-signalChan
	l.Info("Somebody turned off", sig)

}
