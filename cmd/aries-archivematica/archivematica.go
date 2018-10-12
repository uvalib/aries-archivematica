package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
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

func getAIPInfoFromSIPTable(whereClause string) (string, string, error) {
	qs := fmt.Sprintf(`select sipUUID, aipFilename from SIPs where hidden = 0 and %s`, whereClause)

	logger.Printf("query: [%s]", qs)

	rows, err := adb.Query(qs)

	if err != nil {
		logger.Printf("[getAIPInfoFromSIPTable] adb.Query() failed: [%s]", err.Error())
		return "", "", errors.New("AIP from SIP table lookup query failed")
	}

	defer rows.Close()

	var uuids, names []string

	for rows.Next() {
		var sipUUID string
		var aipFilename string
		err = rows.Scan(&sipUUID, &aipFilename)
		if err != nil {
			logger.Printf("[getAIPInfoFromSIPTable] rows.Scan() failed: [%s]", err.Error())
			return "", "", errors.New("AIP from SIP table results scanning failed")
		}
		logger.Printf("sipUUID: [%s]  aipFilename: [%s]", sipUUID, aipFilename)
		uuids = append(uuids, sipUUID)

		// parse out name from aipFilename:
		// format should be 'name.ext' or 'name-sipUUID.ext'
		// remove SIP UUID if present, and drop the extension if any
		name := strings.Replace(aipFilename, "-" + sipUUID, "", 1)
		dot := strings.LastIndex(name,".")
		if dot > 0 {
			name = name[:dot]
		}
		logger.Printf("extracted name: [%s]", name)

		names = append(names, name)
	}

	cnt := len(uuids)

	switch {
	case cnt == 0:
		logger.Printf("[getAIPInfoFromSIPTable] no results")
		return "", "", errors.New("AIP from SIP table query returned no results")
	case cnt > 1:
		logger.Printf("[getAIPInfoFromSIPTable] too many results")
		return "", "", errors.New("AIP from SIP table query returned too many results")
	}

	return uuids[0], names[0], nil
}

func getFileFromUUID(aipUUID string) (string, error) {
	return "fakeFile", nil
}

func getAIPFromId(id string) (*AriesAPI, error) {
	// 1. if id is a UUID, lookup AIP filename in adb; otherwise, lookup UUID
	// 2. extract AIP name from AIP filename: ${AIPName}-${UUID}.7z
	// 2. lookup AIP location in sdb based on UUID

	var aipUUID, aipName, aipFile, where string
	var aipErr error
	var aipInfo AriesAPI

	uuid, uuidErr := uuid.FromString(id)
	if uuidErr != nil {
		// id is not a UUID; lookup by name
		logger.Printf("[%s] is not a UUID; looking up by name...", id)
		where = fmt.Sprintf(`aipFilename regexp '^%s(-[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}\\.|\\.).*$'`, id)
	} else {
		// id is a UUID; lookup by UUID
		logger.Printf("[%s] is a UUID; looking up by UUID [%s]...", id, uuid.String())
		where = fmt.Sprintf(`lcase(sipUUID) = lcase('%s')`, uuid.String())
	}

	aipUUID, aipName, aipErr = getAIPInfoFromSIPTable(where)

	if aipErr != nil {
		logger.Printf("AIP identifier lookup failed: %s", aipErr.Error())
		return nil, errors.New("AIP identifier lookup failed")
	}

	aipFile, aipErr = getFileFromUUID(aipUUID)
	if aipErr != nil {
		logger.Printf("AIP filename lookup failed: %s", aipErr.Error())
		return nil, errors.New("AIP filename lookup failed")
	}

	aipAdminURL := fmt.Sprintf("http://amatica.lib.virginia.edu:81/archival-storage/%s/", aipUUID)

	aipInfo.addIdentifier(aipName)
	aipInfo.addIdentifier(aipUUID)
	aipInfo.addAdministrativeUrl(aipAdminURL)
	aipInfo.addMasterFile(aipFile)

	return &aipInfo, nil
}

/* Handles a request for information about a single ID */
func archivematicaHandleId(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)

	id := params.ByName("id")

//	getAIPsViaApplication()
//	getAIPsViaStorageService()

	archivematicaResponse, err := getAIPFromId(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Printf("AIP not found with ID: %s", id)
		fmt.Fprintf(w, "AIP not found with ID: %s", id)
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
