package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

const progname = "aries-archivematica"
const version = "1.0.0"

var adb, sdb *sql.DB
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

	// Initialize database connections
	var connectStr string
	var err error

	logger.Printf("Initializing Archivematica database connection...")
	connectStr = fmt.Sprintf("%s:%s@%s(%s)/%s", config.applicationDBUser.value, config.applicationDBPass.value,
		config.applicationDBProtocol.value, config.applicationDBHost.value, config.applicationDBName.value)
	adb, err = sql.Open("mysql", connectStr)
	if err != nil {
		fmt.Printf("Archivematica database initialization failed: %s", err.Error())
		os.Exit(1)
	}
	defer adb.Close()

	logger.Printf("Initializing Archivematica Storage Service database connection...")
	connectStr = fmt.Sprintf("%s:%s?mode=ro", config.storageDBProtocol.value, config.storageDBHost.value)
	sdb, err = sql.Open("sqlite3", connectStr)
	if err != nil {
		fmt.Printf("Archivematica Storage Service database initialization failed: %s", err.Error())
		os.Exit(1)
	}
	defer sdb.Close()

	// initialize http client
	client = &http.Client{Timeout: 10 * time.Second}

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", rootHandler)
	mux.GET("/api/aries/:id", archivematicaHandleId)

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
