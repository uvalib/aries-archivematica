package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getAIPs() {
	qs := "select sipUUID, aipFilename from SIPs where hidden = 0 and aipFilename is not null"

	rows, err := db.Query(qs)

	if err != nil {
		logger.Printf("db.Query() failed: [%s]", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var sipUUID string
		var aipFilename string
		err = rows.Scan(&sipUUID, &aipFilename)
		if err != nil {
			logger.Printf("rows.Scan() failed: [%s]", err.Error())
			return
		}
		logger.Printf("sipUUID: [%s]  aipFilename: [%s]", sipUUID, aipFilename)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		logger.Printf("rows.Err() failed: [%s]", err.Error())
		return
	}
}

/* Handles a request for information about a single ID */
func archivematicaHandleId(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)

	id := params.ByName("id")

	getAIPs()

	// build Aries API response object
	var archivematicaResponse AriesAPI
	archivematicaResponse.Identifiers = append(archivematicaResponse.Identifiers, id)

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
