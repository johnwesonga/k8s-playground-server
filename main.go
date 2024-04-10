package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("indexHandler called")
	w.WriteHeader(200)
	w.Write([]byte("Hello World!"))

}

func main() {
	httpListenPort := os.Getenv("PORT")
	if httpListenPort == "" {
		httpListenPort = ":8080"
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/", indexHandler).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	go func() {
		log.Fatal(http.ListenAndServe(httpListenPort, loggedRouter))
	}()

	log.Println("Starting k8s-playground-server on port", httpListenPort)
	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}
	log.Print("The service is shutting down...")
	log.Print("Done")

}
