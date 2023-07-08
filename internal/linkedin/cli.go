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

func (l *LinkedIn) Execute() (data map[string]interface{}, err error) {
	url := fmt.Sprintf("https://www.linkedin.com/in/%s", l.Username)

	if err = l.Get(url); err != nil {
		log.Error().Str("url", url).Err(err).Msg("failed to get the LinkedIn profile")
		return
	}

	data = l.extractUserProfile()
	return
}

func (l *LinkedIn) extractUserProfile() (data map[string]interface{}) {
	data = map[string]interface{}{}

	data["name"] = l.getFieldData(".top-card-layout__title")
	data["title"] = l.getFieldData(".top-card-layout__headline")
	data["location"] = l.getFieldData(".top-card__subline-item")

	// list all the experiences
	selector := ".experience__list > li"
	elms, _ := l.FindElements(selenium.ByCSSSelector, selector)
	for _, elm := range elms {
		// get the company name
		subjectSelector := ".profile-section-card__subtitle > a"
		subjectElm, _ := elm.FindElement(selenium.ByCSSSelector, subjectSelector)

		subject, _ := subjectElm.Text()

		// get the job title
		jobSelector := ".profile-section-card__title"
		jobElm, _ := elm.FindElement(selenium.ByCSSSelector, jobSelector)

		job, _ := jobElm.Text()

		data[subject] = job
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
