// Get the GitHub CLI interface.
package github

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/cmj0121/clotho/internal/utils"
)

// The options for the CLI on Clotho.
type GitHub struct {
	Username string `arg:"" help:"The GitHub username."`

	// The HTTP client for the GitHub API.
	utils.Client
}

// Get the GitHub user information.
func (g *GitHub) Execute() (data map[string]interface{}, err error) {
	var resp *http.Response
	var body []byte

	resp, err = g.Get("https://api.github.com/users/" + g.Username)
	if err != nil {
		log.Warn().Err(err).Str("name", g.Username).Msg("Failed to get the GitHub user.")
		return
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("name", g.Username).Msg("Failed to read response.")
		return
	}

	data = make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Warn().Err(err).Str("name", g.Username).Msg("Failed to parse the GitHub user.")
		return
	}

	return
}
