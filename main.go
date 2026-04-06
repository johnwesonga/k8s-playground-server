package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("indexHandler called")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hello World!")

}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("aboutHandler called")
	w.WriteHeader(200)
	fmt.Fprintf(w, "About Handler")

}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Pass request to the next handler
			next.ServeHTTP(w, r)

			// Log request details
			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("duration", time.Since(start)),
				zap.String("user_agent", r.UserAgent()))
		})
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	httpListenPort := os.Getenv("PORT")
	if httpListenPort == "" {
		httpListenPort = "8080"
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/about", aboutHandler).Methods("GET")

	// loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	loggedRouter := LoggingMiddleware(logger)(router)

	go func() {
		log.Fatal(http.ListenAndServe(":"+httpListenPort, loggedRouter))
	}()

	log.Printf("Starting k8s-playground-server on port %s", httpListenPort)
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
