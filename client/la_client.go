package client

import "net/http"

type LAClient struct {
	httpClient *http.Client
	laHost     string
	laToken    string
}

func NewLAClient(httpClient *http.Client, laHost string, laToken string) *LAClient {
	return &LAClient{
		httpClient: httpClient,
		laHost:     laHost,
		laToken:    laToken,
	}
}
