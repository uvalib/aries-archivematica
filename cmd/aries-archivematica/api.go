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
