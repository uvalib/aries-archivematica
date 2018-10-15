package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/satori/go.uuid"
)

// retrieves and extracts the AIP UUID and name based on UUID or name (via where clause),
// using values in the application database
func getUUIDAndNameFromApplicationDatabase(whereClause string) (string, string, error) {
	qs := fmt.Sprintf(`select sipUUID, aipFilename from SIPs where hidden = 0 and %s`, whereClause)

	rows, err := adb.Query(qs)

	if err != nil {
		logger.Printf("[getUUIDAndNameFromApplicationDatabase] adb.Query() failed: [%s]", err.Error())
		return "", "", errors.New("AIP from SIP table lookup query failed")
	}

	defer rows.Close()

	var uuids, names []string

	for rows.Next() {
		var sipUUID string
		var aipFilename string
		err = rows.Scan(&sipUUID, &aipFilename)
		if err != nil {
			logger.Printf("[getUUIDAndNameFromApplicationDatabase] rows.Scan() failed: [%s]", err.Error())
			return "", "", errors.New("AIP from SIP table results scanning failed")
		}
		logger.Printf("sipUUID: [%s]  aipFilename: [%s]", sipUUID, aipFilename)
		uuids = append(uuids, sipUUID)

		// parse out name from aipFilename:
		// format should be 'name.ext' or 'name-sipUUID.ext'
		// remove SIP UUID if present, and drop the extension if any
		name := strings.Replace(aipFilename, "-"+sipUUID, "", 1)
		dot := strings.LastIndex(name, ".")
		if dot > 0 {
			name = name[:dot]
		}
		logger.Printf("extracted name: [%s]", name)

		names = append(names, name)
	}

	cnt := len(uuids)

	switch {
	case cnt == 0:
		logger.Printf("[getUUIDAndNameFromApplicationDatabase] no results")
		return "", "", errors.New("AIP from SIP table query returned no results")
	case cnt > 1:
		logger.Printf("[getUUIDAndNameFromApplicationDatabase] too many results")
		return "", "", errors.New("AIP from SIP table query returned too many results")
	}

	return uuids[0], names[0], nil
}

// builds the full path to the master file for the given AIP UUID,
// using values in the storage service database
func getMasterFileFromStorageServiceDatabase(aipUUID string) (string, error) {
	qs := fmt.Sprintf(`select s.path || l.relative_path || '/' || p.current_path as file from locations_package p left join locations_location l on p.current_location_id = l.uuid left join locations_space s on l.space_id = s.uuid where l.enabled = 1 and l.purpose = "AS" and p.package_type = "AIP" and p.status = "UPLOADED" and p.uuid = '%s'`, aipUUID)

	rows, err := sdb.Query(qs)

	if err != nil {
		logger.Printf("[getMasterFileFromStorageServiceDatabase] sdb.Query() failed: [%s]", err.Error())
		return "", errors.New("Master file lookup query failed")
	}

	defer rows.Close()

	var files []string

	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err != nil {
			logger.Printf("[getMasterFileFromStorageServiceDatabase] rows.Scan() failed: [%s]", err.Error())
			return "", errors.New("Master file lookup query results scanning failed")
		}
		logger.Printf("file: [%s]", file)

		files = append(files, file)
	}

	cnt := len(files)

	switch {
	case cnt == 0:
		logger.Printf("[getMasterFileFromStorageServiceDatabase] no results")
		return "", errors.New("Master file lookup query returned no results")
	case cnt > 1:
		logger.Printf("[getMasterFileFromStorageServiceDatabase] too many results")
		return "", errors.New("Master file lookup query returned too many results")
	}

	return files[0], nil
}

// takes an ID (name or UUID) and retrieves name/UUID/filename for the matching AIP, if any
func getAIPFromIdViaDatabase(id string) (*AriesAPI, error) {
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

	aipUUID, aipName, aipErr = getUUIDAndNameFromApplicationDatabase(where)

	if aipErr != nil {
		logger.Printf("AIP identifier lookup failed: %s", aipErr.Error())
		return nil, errors.New("AIP identifier lookup failed")
	}

	aipFile, aipErr = getMasterFileFromStorageServiceDatabase(aipUUID)
	if aipErr != nil {
		logger.Printf("AIP filename lookup failed: %s", aipErr.Error())
		return nil, errors.New("AIP filename lookup failed")
	}

	aipAdminURL := strings.Replace(config.adminUrlTemplate.value, "{UUID}", aipUUID, 1)

	aipInfo.addIdentifier(aipName)
	aipInfo.addIdentifier(aipUUID)
	aipInfo.addAdministrativeUrl(aipAdminURL)
	aipInfo.addMasterFile(aipFile)

	return &aipInfo, nil
}