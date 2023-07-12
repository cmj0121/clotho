// Get the Demo CLI interface.
package demo

import (
	"strings"
	"bufio"
	"io"
	"os"
	"fmt"
	"errors"
	"net/http"

	"github.com/cmj0121/clotho/internal/utils"
	"github.com/rs/zerolog/log"
)

type Demo struct {
	Link   string `arg:"" help:"The target link."`
	Action string `enum:"http,selenium" required:"" default:"http" help:"The action to perform."`

	// The HTTP client for the GitHub API.
	utils.Client
	// The Selenium wrapper client.
	utils.Selenium
}

func (d *Demo) Prologue() {
	switch d.Action {
	case "http":
		d.Client.Prologue()
	case "selenium":
		d.Selenium.Prologue()
	}
}

func (d *Demo) Epilogue() {
	switch d.Action {
	case "http":
		d.Client.Epilogue()
	case "selenium":
		d.Selenium.Epilogue()
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
	case "selenium":
		help := strings.Join([]string{
			"exit: exit the demo command",
		}, "\n")

		reader := bufio.NewReader(os.Stdin)
		d.Selenium.Get(d.Link)
		for {
			fmt.Print(">>> ")

			text, rerr := reader.ReadString('\n')
			switch rerr {
				case nil:
					text = strings.TrimSpace(text)
					switch text {
					case "exit":
						return
					default:
						fmt.Printf("unknown command: %v\n", text)
						fmt.Println(help)
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
