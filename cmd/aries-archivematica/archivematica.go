package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func getAIPsViaApplication() {
	qs := "select sipUUID, aipFilename from SIPs where hidden = 0 and aipFilename is not null"

	rows, err := adb.Query(qs)

	if err != nil {
		logger.Printf("adb.Query() failed: [%s]", err.Error())
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

func getAIPsViaStorageService() {
	qs := `select p.uuid, s.path || l.relative_path || '/' || p.current_path as fullpath from locations_package p left join locations_location l on p.current_location_id = l.uuid left join locations_space s on l.space_id = s.uuid where l.enabled = 1 and l.purpose = "AS" and p.package_type = "AIP" and p.status = "UPLOADED"`

	rows, err := sdb.Query(qs)

	if err != nil {
		logger.Printf("sdb.Query() failed: [%s]", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var uuid string
		var path string
		err = rows.Scan(&uuid, &path)
		if err != nil {
			logger.Printf("rows.Scan() failed: [%s]", err.Error())
			return
		}
		logger.Printf("uuid: [%s]  path: [%s]", uuid, path)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		logger.Printf("rows.Err() failed: [%s]", err.Error())
		return
	}
}

func getAIPInfo(id string) (string, string, string, error) {
	// 1. if id is a UUID, lookup AIP filename in adb; otherwise, lookup UUID
	// 2. extract AIP name from AIP filename: ${AIPName}-${UUID}.7z
	// 2. lookup AIP location in sdb based on UUID

	return "00000000-0000-0000-0000-000000000000", "fakeAIP", "/path/to/fakeAIP.7z", nil
}

/* Handles a request for information about a single ID */
func archivematicaHandleId(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)

	id := params.ByName("id")

	getAIPsViaApplication()
	getAIPsViaStorageService()

	aipUUID, aipName, aipFile, aipErr := getAIPInfo(id)
	if aipErr != nil {
		logger.Printf("aipErr: [%s]", aipErr.Error())
	}

	logger.Printf("aipUUID: [%s]  aipName: [%s]  aipFile: [%s]", aipUUID, aipName, aipFile)

	// build Aries API response object
	var archivematicaResponse AriesAPI

	archivematicaResponse.addIdentifier(id)

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
