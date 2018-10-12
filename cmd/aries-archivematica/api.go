package main

type ServiceUrl struct {
	Url      string `json:"url,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

type MetadataUrl struct {
	Url    string `json:"url,omitempty"`
	Schema string `json:"schema,omitempty"`
}

type AriesAPI struct {
	Identifiers        []string      `json:"identifier"`
	AdministrativeUrls []string      `json:"administrative_url,omitempty"`
	AccessUrls         []string      `json:"access_url,omitempty"`
	ServiceUrls        []ServiceUrl  `json:"service_url,omitempty"`
	MetadataUrls       []MetadataUrl `json:"metadata_url,omitempty"`
	MasterFiles        []string      `json:"master_file,omitempty"`
	DerivativeFiles    []string      `json:"derivative_file,omitempty"`
	AccessRestriction  string        `json:"access_restriction,omitempty"`
}

func (a *AriesAPI) addIdentifier(s string) {
	a.Identifiers = append(a.Identifiers, s)
}

func (a *AriesAPI) addAdministrativeUrl(s string) {
	a.AdministrativeUrls = append(a.AdministrativeUrls, s)
}

func (a *AriesAPI) addAccessUrl(s string) {
	a.AccessUrls = append(a.AccessUrls, s)
}

func (a *AriesAPI) addServiceUrl(s ServiceUrl) {
	a.ServiceUrls = append(a.ServiceUrls, s)
}

func (a *AriesAPI) addMetadataUrl(m MetadataUrl) {
	a.MetadataUrls = append(a.MetadataUrls, m)
}

func (a *AriesAPI) addMasterFile(s string) {
	a.MasterFiles = append(a.MasterFiles, s)
}

func (a *AriesAPI) addDerivativeFile(s string) {
	a.DerivativeFiles = append(a.DerivativeFiles, s)
}

func (a *AriesAPI) setAccessRestriction(s string) {
	a.AccessRestriction = s
}
