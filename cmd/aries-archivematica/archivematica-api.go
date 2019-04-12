package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/satori/go.uuid"
)

// storage service response structs
// note: more fields exist, we just restrict to the fields we are interested in
type storageServiceMeta struct {
	TotalCount int `json:"total_count,omitempty"`
}

type storageServiceObject struct {
	CurrentFullPath string `json:"current_full_path,omitempty"`
	UUID            string `json:"uuid,omitempty"`
}

type storageServiceResponse struct {
	Meta    storageServiceMeta     `json:"meta,omitempty"`
	Objects []storageServiceObject `json:"objects,omitempty"`
}

// retrieves the full path to the master file for the given AIP UUID
func getMasterFileFromStorageServiceAPI(aipUUID string) (string, string, error) {

	url := strings.Replace(config.storageAPIUrlTemplate.value, "{UUID}", aipUUID, 1)

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		logger.Printf("NewRequest() failed: %s", reqErr.Error())
		return "", "", errors.New("Failed to create new request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s:%s", config.storageAPIUser.value, config.storageAPIKey.value))

	res, resErr := client.Do(req)
	if resErr != nil {
		logger.Printf("client.Do() failed: %s", resErr.Error())
		return "", "", errors.New("Failed to read response")
	}

	defer res.Body.Close()

	// parse json from body

	var ssResp storageServiceResponse

	decoder := json.NewDecoder(res.Body)
	decErr := decoder.Decode(&ssResp)
	if decErr != nil {
		logger.Printf("Decode() failed: %s", decErr.Error())
		return "", "", errors.New("Failed to decode response")
	}

	// ensure just one result
	switch ssResp.Meta.TotalCount {
	case 0:
		logger.Printf("TotalCount %d / len() %d [no results]", ssResp.Meta.TotalCount, len(ssResp.Objects))
		return "", "", errors.New("No AIP found with this ID")
	case 1:
		logger.Printf("TotalCount %d / len() %d [ok]", ssResp.Meta.TotalCount, len(ssResp.Objects))
		return ssResp.Objects[0].UUID, ssResp.Objects[0].CurrentFullPath, nil
	default:
		logger.Printf("TotalCount %d / len() %d [too many results]", ssResp.Meta.TotalCount, len(ssResp.Objects))
		return "", "", errors.New("Too many AIPs found with this ID")
	}
}

// takes an ID (UUID) and retrieves UUID/filename for the matching AIP, if any
func getAIPFromIdViaAPI(id string) (*AriesAPI, error) {
	var aipInfo AriesAPI

	uuid, uuidErr := uuid.FromString(id)
	if uuidErr != nil {
		logger.Printf("UUID parsing failed: %s", uuidErr.Error())
		return nil, errors.New("Identifier is not a valid UUID")
	}

	aipUUID, aipFile, aipErr := getMasterFileFromStorageServiceAPI(uuid.String())
	if aipErr != nil {
		logger.Printf("AIP filename lookup failed: %s", aipErr.Error())
		return nil, aipErr
	}

	aipAdminURL := strings.Replace(config.adminUrlTemplate.value, "{UUID}", aipUUID, 1)

	aipInfo.addIdentifier(aipUUID)
	aipInfo.addAdministrativeUrl(aipAdminURL)
	aipInfo.addMasterFile(aipFile)

	return &aipInfo, nil
}
