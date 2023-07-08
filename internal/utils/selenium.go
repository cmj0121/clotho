// Get the global utility for the internal packages.
package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tebeka/selenium"
)

type Selenium struct {
	// The path to the ChromeDriver.
	Driver string `help:"The path to the ChromeDriver." default:"chromedriver"`
	// The bind port for the ChromeDriver.
	Port int `help:"The bind address for the ChromeDriver." default:"9515"`

	// The Selenium service for the ChromeDriver.
	*selenium.Service `kong:"-"`

	// The Selenium web driver for the ChromeDriver.
	selenium.WebDriver `kong:"-"`
}

// open the necessary resources for the HTTP client.
func (s *Selenium) Prologue() {
	opts := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(s.Driver, s.Port, opts...)
	if err != nil {
		log.Error().Err(err).Msg("failed to start the ChromeDriver service")
		panic(err)
	}

	s.Service = service

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", s.Port))
	if err != nil {
		log.Error().Err(err).Msg("failed to start the ChromeDriver")
		panic(err)
	}

	s.WebDriver = wd
}

// clean up the resources for the HTTP client.
func (s *Selenium) Epilogue() {
	if s.WebDriver != nil {
		s.WebDriver.Quit()
		s.WebDriver = nil
	}

	if s.Service != nil {
		s.Service.Stop()
		s.Service = nil
	}
}
