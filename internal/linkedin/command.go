// Get the LinkedIn CLI interface.
package linkedin

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmj0121/clotho/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/tebeka/selenium"
)

type LinkedIn struct {
	Username string        `arg:"" help:"The GitHub username."`
	Wait     time.Duration `help:"The wait time for the Selenium dom visible." default:"2s"`

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
		{"name", l.getFieldData(".top-card-layout__title"), "", ""},
		{"title", l.getFieldData(".top-card-layout__headline"), "", ""},
		{"location", l.getFieldData(".top-card__subline-item"), "", ""},
	}

	// wait untils the experience section is loaded
	selector := ".experience__list > li"
	exists := l.WaitVisible(selector, l.Wait)
	if !exists {
		log.Warn().Str("wait", l.Wait.String()).Msg("failed to find the experience section")
		return
	}

	// get the experience section
	elms, _ := l.FindElements(selenium.ByCSSSelector, selector)
	for _, elm := range elms {
		class, err := elm.GetAttribute("class")
		if err != nil {
			log.Warn().Err(err).Msg("failed to get the class attribute")
			continue
		}

		if strings.Contains(class, "profile-section-card") {
			company, _ := l.getText(elm, ".profile-section-card__subtitle > a", "h4.profile-section-card__subtitle")
			title, _ := l.getText(elm, "h3.profile-section-card__title")

			time_ranges_dom, _ := elm.FindElements(selenium.ByCSSSelector, ".date-range > time")
			time_ranges := []string{}
			for _, time_range_dom := range time_ranges_dom {
				time_range, _ := time_range_dom.GetAttribute("innerText")
				time_ranges = append(time_ranges, time_range)
			}

			data = append(data, []string{"", company, title, strings.Join(time_ranges, " - ")})
			continue
		}

		if strings.Contains(class, "experience-group") {
			company, _ := l.getText(elm, "h4.experience-group-header__company")

			selector := ".experience-group-position"
			positions, _ := l.FindElements(selenium.ByCSSSelector, selector)
			for _, position := range positions {
				title, _ := l.getText(position, "h3.profile-section-card__title")

				time_ranges_dom, _ := position.FindElements(selenium.ByCSSSelector, ".date-range > time")
				time_ranges := []string{}
				for _, time_range_dom := range time_ranges_dom {
					time_range, _ := time_range_dom.GetAttribute("innerText")
					time_ranges = append(time_ranges, time_range)
				}

				data = append(data, []string{"", company, title, strings.Join(time_ranges, " - ")})
			}

			continue
		}

		log.Info().Str("class", class).Msg("unknown element")
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

func (l *LinkedIn) getText(elm selenium.WebElement, selectors ...string) (value string, err error) {
	var dom selenium.WebElement

	for _, selector := range selectors {
		dom, err = elm.FindElement(selenium.ByCSSSelector, selector)
		if err != nil {
			log.Debug().Str("selector", selector).Err(err).Msg("failed to find the element")
			continue
		}

		value, err = dom.GetAttribute("innerText")
		if err != nil {
			log.Error().Err(err).Msg("failed to get the text of the element")
			continue
		}

		err = nil
		break
	}

	return
}
