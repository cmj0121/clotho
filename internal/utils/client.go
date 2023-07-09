// Get the global utility for the internal packages.
package utils

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type Client struct {
	// The HTTP client for the GitHub API.
	client *http.Client
}

// open the necessary resources for the HTTP client.
func (c *Client) Prologue() {
	c.client = &http.Client{}
}

// clean up the resources for the HTTP client.
func (c *Client) Epilogue() {
	c.client = nil
}

// send the HTTP GET request to the target URL.
func (c *Client) Get(url string) (resp *http.Response, err error) {
	resp, err = c.client.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("failed to send the HTTP request")
		return
	}

	return
}
