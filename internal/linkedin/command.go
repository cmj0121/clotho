// Get the LinkedIn CLI interface.
package linkedin

import (
	"fmt"

	"github.com/cmj0121/clotho/internal/utils"
	"github.com/tebeka/selenium"
	"github.com/rs/zerolog/log"
)

type LinkedIn struct {
	Username string `arg:"" help:"The GitHub username."`

	// The Selenium wrapper client.
	utils.Selenium
}

func (l *LinkedIn) Execute() (resp interface{}, err error) {
	url := fmt.Sprintf("https://www.linkedin.com/in/%s", l.Username)

	if err = l.Get(url); err != nil {
		log.Error().Str("url", url).Err(err).Msg("failed to get the LinkedIn profile")
		return
	}

	resp = l.extractUserProfile()
	return
}

func (l *LinkedIn) extractUserProfile() (data [][]string) {
	data = [][]string{
		[]string{"name", l.getFieldData(".top-card-layout__title")},
		[]string{"title", l.getFieldData(".top-card-layout__headline")},
		[]string{"location", l.getFieldData(".top-card__subline-item")},
	}

	return
}

func (l *LinkedIn) getFieldData(selector string) (value string) {
	element, err := l.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Info().Str("selector", selector).Err(err).Msg("failed to find the element")
		return
	}

	value, err = element.Text()
	if err != nil {
		log.Error().Str("selector", selector).Err(err).Msg("failed to get the text of the element")
		return
	}

	return
}
