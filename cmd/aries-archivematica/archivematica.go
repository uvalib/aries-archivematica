package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getAIPFromId(id string) (*AriesAPI, error) {
	return getAIPFromIdViaAPI(id)
}

/* Handles a request for information about a single ID */
func archivematicaIdHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)

	id := params.ByName("id")

	archivematicaResponse, err := getAIPFromId(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Printf("%s", err.Error())
		fmt.Fprintf(w, "%s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, jerr := json.Marshal(archivematicaResponse)
	if jerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("JSON marshal failed: [%s]", jerr.Error())
		fmt.Fprintf(w, "JSON marshal failed")
		return
	}

	fmt.Fprintf(w, string(j))
}
