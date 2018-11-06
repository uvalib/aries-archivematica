package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const progname = "aries-archivematica"
const version = "1.0.0"

var logger *log.Logger
var client *http.Client

/**
 * Main entry point for the web service
 */
func main() {
	// use below to log to console....
	logger = log.New(os.Stdout, "", log.LstdFlags)

	// Load cfg
	logger.Printf("===> %s %s staring up <===", progname, version)
	logger.Printf("Loadiing configuration...")
	getConfigValues()

	// initialize http client
	client = &http.Client{Timeout: 10 * time.Second}

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", rootHandler)
	mux.GET("/api/aries", apiHandler)
	mux.GET("/api/aries/:id", archivematicaIdHandler)

	logger.Printf("Start service on port %s", config.listenPort.value)

	if config.useHttps.value == true {
		log.Fatal(http.ListenAndServeTLS(":"+config.listenPort.value, config.sslCrt.value, config.sslKey.value, cors.Default().Handler(mux)))
	} else {
		log.Fatal(http.ListenAndServe(":"+config.listenPort.value, cors.Default().Handler(mux)))
	}
}

// Handle a request for /
func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)
	fmt.Fprintf(w, "%s version %s", progname, version)
}

// Handle a request for /api/aries
func apiHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)
	fmt.Fprintf(w, "Archivematica Aries API")
}
