package repository

import "net/http"

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(
	baseURL string,
	httpClient *http.Client,
) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}
