package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const progname = "aries-archivematica"
const version = "0.0.1"

var logger *log.Logger

/**
 * Main entry point for the web service
 */
func main() {
	// use below to log to console....
	logger = log.New(os.Stdout, "", log.LstdFlags)

	// Load cfg
	logger.Printf("===> %s staring up <===", progname)
	logger.Printf("Load configuration...")
	getConfigValues()

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", rootHandler)
	mux.GET("/archivematica/:id", archivematicaHandleId)

	logger.Printf("Start service on port %s", config.listenPort.value)

	if config.useHttps.value == true {
		log.Fatal(http.ListenAndServeTLS(":"+config.listenPort.value, config.sslCrt.value, config.sslKey.value, cors.Default().Handler(mux)))
	} else {
		log.Fatal(http.ListenAndServe(":"+config.listenPort.value, cors.Default().Handler(mux)))
	}
}

/**
 * Handle a request for /
 */
func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)
	fmt.Fprintf(w, "%s version %s", progname, version)
}
