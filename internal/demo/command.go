// Get the Demo CLI interface.
package demo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cmj0121/clotho/internal/utils"
	"github.com/rs/zerolog/log"
)

type Demo struct {
	Link   string `arg:"" help:"The target link."`
	Action string `enum:"http,chrome" required:"" default:"http" help:"The action to perform."`

	// The HTTP client for the GitHub API.
	utils.Client
	// The Chrome wrapper client.
	utils.Chrome
}

func (d *Demo) Prologue() {
	switch d.Action {
	case "http":
		d.Client.Prologue()
	case "chrome":
		d.Chrome.Prologue()
	}
}

func (d *Demo) Epilogue() {
	switch d.Action {
	case "http":
		d.Client.Epilogue()
	case "chrome":
		d.Chrome.Epilogue()
	}
}

// Get the GitHub user information.
func (d *Demo) Execute() (data interface{}, err error) {
	log.Info().Str("action", d.Action).Msg("execute the demo command")

	switch d.Action {
	case "http":
		var resp *http.Response

		resp, err = d.Client.Get(d.Link)
		if err != nil {
			log.Error().Err(err).Msg("failed to get the link")
			return
		}
		defer resp.Body.Close()

		var body []byte
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read the response body")
			return
		}

		fmt.Println(string(body))
	case "chrome":
		reader := bufio.NewReader(os.Stdin)
		link := d.Link
		for {
			d.Chrome.Navigate(link)

			fmt.Print(">>> ")
			text, rerr := reader.ReadString('\n')

			switch rerr {
			case nil:
				text = strings.TrimSpace(text)

				switch {
				case strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://"):
					link = text
				default:
					fmt.Printf("invalid link: %s\n", text)
				}
			case io.EOF:
				return
			default:
				log.Error().Err(rerr).Msg("failed to read the command")

				err = errors.New(fmt.Sprintf("failed to read the command: %v", rerr))
				return
			}
		}

	default:
		err = errors.New(fmt.Sprintf("unknown action: %v", d.Action))
		return
	}

	return
}
